package config

import "fnchatbot/internal/models"

type ProviderDefinition struct {
	ProviderID string                    `json:"provider_id"`
	Name       string                    `json:"name"`
	Type       models.ProviderType       `json:"type"`
	BaseURL    string                    `json:"base_url"`
	ApiOptions models.ProviderApiOptions `json:"api_options"`
}

var SystemProviders = []ProviderDefinition{
	{ProviderID: "openai", Name: "OpenAI", Type: models.ProviderTypeOpenAIResponse, BaseURL: "https://api.openai.com"},
	{ProviderID: "anthropic", Name: "Anthropic", Type: models.ProviderTypeAnthropic, BaseURL: "https://api.anthropic.com"},
	{ProviderID: "gemini", Name: "Google Gemini", Type: models.ProviderTypeGemini, BaseURL: "https://generativelanguage.googleapis.com"},
	{ProviderID: "grok", Name: "Grok", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.x.ai"},
	{ProviderID: "mistral", Name: "Mistral", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.mistral.ai"},
	{ProviderID: "groq", Name: "Groq", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.groq.com/openai"},
	{ProviderID: "perplexity", Name: "Perplexity", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.perplexity.ai"},
	{ProviderID: "nvidia", Name: "NVIDIA", Type: models.ProviderTypeOpenAI, BaseURL: "https://integrate.api.nvidia.com"},
	{ProviderID: "together", Name: "Together", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.together.xyz"},
	{ProviderID: "fireworks", Name: "Fireworks", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.fireworks.ai/inference"},
	{ProviderID: "hyperbolic", Name: "Hyperbolic", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.hyperbolic.xyz"},
	{ProviderID: "cerebras", Name: "Cerebras", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.cerebras.ai/v1"},

	{ProviderID: "zhipu", Name: "智谱 AI", Type: models.ProviderTypeOpenAI, BaseURL: "https://open.bigmodel.cn/api/paas/v4/"},
	{ProviderID: "deepseek", Name: "DeepSeek", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.deepseek.com"},
	{ProviderID: "moonshot", Name: "Moonshot AI", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.moonshot.cn"},
	{ProviderID: "baichuan", Name: "百川 AI", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.baichuan-ai.com"},
	{ProviderID: "dashscope", Name: "阿里百炼", Type: models.ProviderTypeOpenAI, BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1/"},
	{ProviderID: "doubao", Name: "豆包", Type: models.ProviderTypeOpenAI, BaseURL: "https://ark.cn-beijing.volces.com/api/v3/"},
	{ProviderID: "minimax", Name: "MiniMax", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.minimaxi.com/v1/"},
	{ProviderID: "hunyuan", Name: "腾讯混元", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.hunyuan.cloud.tencent.com"},
	{ProviderID: "baidu-cloud", Name: "百度云", Type: models.ProviderTypeOpenAI, BaseURL: "https://qianfan.baidubce.com/v2/"},
	{ProviderID: "yi", Name: "零一万物", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.lingyiwanwu.com"},
	{ProviderID: "stepfun", Name: "阶跃星辰", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.stepfun.com"},
	{ProviderID: "mimo", Name: "小米 MiMo", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.xiaomimimo.com"},

	{ProviderID: "aihubmix", Name: "AiHubMix", Type: models.ProviderTypeOpenAI, BaseURL: "https://aihubmix.com"},
	{ProviderID: "silicon", Name: "Silicon", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.siliconflow.cn"},
	{ProviderID: "openrouter", Name: "OpenRouter", Type: models.ProviderTypeOpenAI, BaseURL: "https://openrouter.ai/api/v1/"},
	{ProviderID: "302ai", Name: "302.AI", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.302.ai"},
	{ProviderID: "tokenflux", Name: "TokenFlux", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.tokenflux.ai/openai/v1"},

	{ProviderID: "ollama", Name: "Ollama", Type: models.ProviderTypeOllama, BaseURL: "http://localhost:11434"},
	{ProviderID: "lmstudio", Name: "LM Studio", Type: models.ProviderTypeOpenAI, BaseURL: "http://localhost:1234"},
	{ProviderID: "gpustack", Name: "GPUStack", Type: models.ProviderTypeOpenAI, BaseURL: ""},
	{ProviderID: "ovms", Name: "OpenVINO Model Server", Type: models.ProviderTypeOpenAI, BaseURL: "http://localhost:8000/v3/"},

	{ProviderID: "azure-openai", Name: "Azure OpenAI", Type: models.ProviderTypeAzureOpenAI, BaseURL: ""},
	{ProviderID: "aws-bedrock", Name: "AWS Bedrock", Type: models.ProviderTypeAWBedrock, BaseURL: ""},
	{ProviderID: "vertexai", Name: "Vertex AI", Type: models.ProviderTypeVertexAI, BaseURL: ""},
	{ProviderID: "github", Name: "GitHub Models", Type: models.ProviderTypeOpenAI, BaseURL: "https://models.github.ai/inference"},
	{ProviderID: "copilot", Name: "GitHub Copilot", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.githubcopilot.com/"},

	{ProviderID: "jina", Name: "Jina", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.jina.ai"},
	{ProviderID: "voyageai", Name: "VoyageAI", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.voyageai.com"},
	{ProviderID: "huggingface", Name: "Hugging Face", Type: models.ProviderTypeOpenAIResponse, BaseURL: "https://router.huggingface.co/v1/"},
	{ProviderID: "gateway", Name: "Vercel AI Gateway", Type: models.ProviderTypeGateway, BaseURL: "https://ai-gateway.vercel.sh/v1/ai"},

	{ProviderID: "ocoolai", Name: "ocoolAI", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.ocoolai.com"},
	{ProviderID: "ppio", Name: "PPIO", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.ppinfra.com/v3/openai/"},
	{ProviderID: "alayanew", Name: "AlayaNew", Type: models.ProviderTypeOpenAI, BaseURL: "https://deepseek.alayanew.com"},
	{ProviderID: "qiniu", Name: "Qiniu", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.qnaigc.com"},
	{ProviderID: "dmxapi", Name: "DMXAPI", Type: models.ProviderTypeOpenAI, BaseURL: "https://www.dmxapi.cn"},
	{ProviderID: "burncloud", Name: "BurnCloud", Type: models.ProviderTypeOpenAI, BaseURL: "https://ai.burncloud.com"},
	{ProviderID: "cephalon", Name: "Cephalon", Type: models.ProviderTypeOpenAI, BaseURL: "https://cephalon.cloud/user-center/v1/model"},
	{ProviderID: "lanyun", Name: "LANYUN", Type: models.ProviderTypeOpenAI, BaseURL: "https://maas-api.lanyun.net"},
	{ProviderID: "ph8", Name: "PH8", Type: models.ProviderTypeOpenAI, BaseURL: "https://ph8.co"},
	{ProviderID: "sophnet", Name: "SophNet", Type: models.ProviderTypeOpenAI, BaseURL: "https://www.sophnet.com/api/open-apis/v1"},
	{ProviderID: "aionly", Name: "AIOnly", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.aiionly.com"},
	{ProviderID: "longcat", Name: "LongCat", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.longcat.chat/openai"},
	{ProviderID: "infini", Name: "Infini", Type: models.ProviderTypeOpenAI, BaseURL: "https://cloud.infini-ai.com/maas"},
	{ProviderID: "modelscope", Name: "ModelScope", Type: models.ProviderTypeOpenAI, BaseURL: "https://api-inference.modelscope.cn/v1/"},
	{ProviderID: "xirang", Name: "Xirang", Type: models.ProviderTypeOpenAI, BaseURL: "https://wishub-x1.ctyun.cn"},
	{ProviderID: "tencent-cloud-ti", Name: "Tencent Cloud TI", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.lkeap.cloud.tencent.com"},
	{ProviderID: "poe", Name: "Poe", Type: models.ProviderTypeOpenAI, BaseURL: "https://api.poe.com/v1/"},
	{ProviderID: "new-api", Name: "New API", Type: models.ProviderTypeNewAPI, BaseURL: "http://localhost:3000"},
}

func GetProviderByID(providerID string) *ProviderDefinition {
	for i := range SystemProviders {
		if SystemProviders[i].ProviderID == providerID {
			return &SystemProviders[i]
		}
	}
	return nil
}

func GetProvidersByType(providerType models.ProviderType) []ProviderDefinition {
	var result []ProviderDefinition
	for i := range SystemProviders {
		if SystemProviders[i].Type == providerType {
			result = append(result, SystemProviders[i])
		}
	}
	return result
}
