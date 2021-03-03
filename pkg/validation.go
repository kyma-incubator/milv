package pkg

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/schollz/closestmatch"
)

type Waiter interface {
	Wait()
}

type Validator struct {
	client http.Client
	waiter Waiter
}

func NewValidator(client http.Client, limiter Waiter) *Validator {
	return &Validator{client: client, waiter: limiter}
}

func (v *Validator) Links(links []Link, optionalHeaders ...Headers) []Link {
	if len(links) == 0 {
		return []Link{}
	}

	var headers Headers
	var headersExist bool
	if len(optionalHeaders) == 1 {
		headers = optionalHeaders[0]
		headersExist = len(headers) > 0
	}

	var validatedLinks []Link
	for _, link := range links {
		if link.TypeOf == ExternalLink {
			link, _ = v.externalLink(link)
			validatedLinks = append(validatedLinks, link)
		} else if link.TypeOf == InternalLink {
			link, _ = v.internalLink(link)
			validatedLinks = append(validatedLinks, link)
		} else {
			if headersExist {
				link, _ = v.hashInternalLink(link, headers)
				validatedLinks = append(validatedLinks, link)
			}
		}
	}
	return validatedLinks
}

func (v *Validator) ExternalLinks(links Links) (Links, error) {
	for _, link := range links {
		if link.TypeOf == ExternalLink {
			link, _ = v.externalLink(link)
		}
	}
	return links, nil
}

func (v *Validator) InternalLinks(links Links) (Links, error) {
	for _, link := range links {
		if link.TypeOf == InternalLink {
			link, _ = v.externalLink(link)
		}
	}
	return links, nil
}

func (v *Validator) HashInternalLinks(links Links, headers Headers) (Links, error) {
	for _, link := range links {
		if link.TypeOf == HashInternalLink {
			link, _ = v.hashInternalLink(link, headers)
		}
	}
	return links, nil
}

func (v *Validator) externalLink(link Link) (Link, error) {
	if link.TypeOf != ExternalLink {
		return link, nil
	}

	var status bool
	message := ""

	url, err := url.Parse(link.AbsPath)
	if err != nil {
		link.Result.Status = false
		link.Result.Message = err.Error()
		return link, err
	}
	absPath := fmt.Sprintf("%s://%s%s", url.Scheme, url.Host, url.Path)

	if link.Config != nil && link.Config.Timeout != nil && *link.Config.Timeout != 0 {
		v.client.Timeout = time.Duration(int(time.Second) * (*link.Config.Timeout))
	} else {
		v.client.Timeout = time.Duration(int(time.Second) * 30)
	}

	requestRepeats := 1
	if link.Config != nil && link.Config.RequestRepeats != nil && *link.Config.RequestRepeats > 0 {
		requestRepeats = *link.Config.RequestRepeats
	}

	for i := 0; i < requestRepeats; i++ {
		resp, err := v.client.Get(absPath)
		if err != nil {
			status = false
			message = err.Error()
			continue
		}

		allowRedirect := false
		if link.Config != nil && link.Config.AllowRedirect != nil {
			allowRedirect = *link.Config.AllowRedirect
		}

		statusCode, http2xxPattern := strconv.Itoa(resp.StatusCode), `^2[0-9][0-9]`
		if allowRedirect {
			http2xxPattern = `^2[0-9][0-9]|^3[0-9][0-9]`
		}

		if match, _ := regexp.MatchString(http2xxPattern, statusCode); match && resp != nil {
			status = true

			if !allowRedirect && url.Fragment != "" {
				match, _ = regexp.MatchString(`[a-zA-Z]`, string(url.Fragment[0]))
				if !match {
					break
				}

				parser := &Parser{}
				anchors := parser.Anchors(resp.Body)

				if contains(anchors, url.Fragment) {
					status = true
				} else {
					cm := closestmatch.New(anchors, []int{4, 3, 5})
					closestAnchor := cm.Closest(url.Fragment)

					status = false
					if closestAnchor != "" {
						message = fmt.Sprintf("The specified anchor doesn't exist in website. Did you mean about #%s?", closestAnchor)
					} else {
						message = "The specified anchor doesn't exist"
					}
				}
			}

			CloseBody(resp.Body)
			break
		} else if resp.StatusCode == http.StatusTooManyRequests {
			status = false
			message = "Too many requests"
			v.waiter.Wait()
			CloseBody(resp.Body)
			continue
		} else {
			status = false
			message = resp.Status
			CloseBody(resp.Body)
		}
	}

	link.Result.Status = status
	link.Result.Message = message
	return link, nil
}

func (v *Validator) internalLink(link Link) (Link, error) {
	if link.TypeOf != InternalLink {
		return link, nil
	}

	splitted := strings.Split(link.AbsPath, "#")

	if err := fileExists(splitted[0]); err == nil {
		link.Result.Status = true

		if len(splitted) == 2 {
			if !v.isHashInFile(splitted[0], splitted[1]) {
				link.Result.Status = false
				link.Result.Message = "The specified header doesn't exist in file"
			}
		}
	} else {
		link.Result.Status = false
		link.Result.Message = "The specified file doesn't exist"
	}
	return link, nil
}

func (*Validator) hashInternalLink(link Link, headers Headers) (Link, error) {
	if link.TypeOf != HashInternalLink {
		return link, nil
	}

	if match := headerExists(link.RelPath, headers); match {
		link.Result.Status = true
	} else {
		link.Result.Status = false
		link.Result.Message = "The specified header doesn't exist in file"
	}
	return link, nil
}

func (*Validator) isHashInFile(file, header string) bool {
	markdown, err := readMarkdown(file)
	if err != nil {
		return false
	}

	parser := Parser{}
	return headerExists(header, parser.Headers(markdown))
}
