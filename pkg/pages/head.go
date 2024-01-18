package pages

import "mrdmitry/blog/pkg/monke"

type HeadSnippet struct {
	GitRev string
}

func NewHeadSnippet() HeadSnippet {
	revision, _ := monke.GitRevision()
	return HeadSnippet{
		GitRev: revision,
	}
}
