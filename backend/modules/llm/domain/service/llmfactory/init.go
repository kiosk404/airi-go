package llmfactory

func NewFactory() IFactory {
	return &FactoryImpl{}
}
