package lang

type Token int

type Keyword struct {
	Token         Token
	FutureKeyword bool
	Strict        bool
}

type KeywordMap map[string]*Keyword

var langKeyWordMap map[string]KeywordMap = make(map[string]KeywordMap)

func LangKeyWordMap(lang string) (kwrdmap KeywordMap) {
	if lang != "" {
		kwrdmap, _ = langKeyWordMap[lang]
	}
	return
}

func RegisterLangKeyWordMap(lang string, kwrdmap ...KeywordMap) {
	if lang != "" {
		var crntkwdmp KeywordMap = LangKeyWordMap(lang)
		var addlwddmp = crntkwdmp == nil
		for _, kwrdmp := range kwrdmap {
			if crntkwdmp == nil {
				crntkwdmp = kwrdmp
			} else {
				for kwrdnm, kwrd := range kwrdmp {
					if _, kwrdok := crntkwdmp[kwrdnm]; !kwrdok {
						crntkwdmp[kwrdnm] = kwrd
					} else {

					}
				}
			}
		}
		if addlwddmp {
			langKeyWordMap[lang] = crntkwdmp
		}
	}
}
