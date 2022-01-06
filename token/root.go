package token

type rootToken struct {
}

func (token *rootToken) Apply(root, current interface{}, next []Token) (interface{}, error) {
	if len(next) > 0 {
		return next[0].Apply(root, root, next[1:])
	}
	return root, nil
}
