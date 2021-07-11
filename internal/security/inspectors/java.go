package inspectors

type JavaHook struct{}

func (hook *JavaHook) After(filePath string, line string, secret string) {}

func (hook *JavaHook) Before() {}
