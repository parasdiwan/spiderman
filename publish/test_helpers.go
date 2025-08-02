package publish

type TestPublisher struct {
	Published []string
}

func NewTestPublisher() *TestPublisher {
	return &TestPublisher{
		Published: make([]string, 0),
	}
}

func (p *TestPublisher) Publish(string string, strings []string) error {
	p.Published = append(p.Published, string)
	for _, str := range strings {
		p.Published = append(p.Published, str)
	}
	return nil
}

func (c *TestPublisher) RecordError(_ string, _ ErrType) error {
	return nil
}

func (p *TestPublisher) PublishStats() error {
	return nil
}
