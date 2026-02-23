export interface PredefinedProvider {
  id: string
  name: string
  type: string
  baseUrl: string
  icon?: string
  description?: string
}

export const predefinedProviders: PredefinedProvider[] = [
  { id: 'openai', name: 'OpenAI', type: 'openai', baseUrl: 'https://api.openai.com/v1', description: 'GPT-4, GPT-3.5, DALL-E' },
  { id: 'anthropic', name: 'Anthropic', type: 'anthropic', baseUrl: 'https://api.anthropic.com', description: 'Claude series models' },
  { id: 'google', name: 'Google AI', type: 'gemini', baseUrl: 'https://generativelanguage.googleapis.com', description: 'Gemini, PaLM' },
  { id: 'azure-openai', name: 'Azure OpenAI', type: 'azure-openai', baseUrl: 'https://YOUR_RESOURCE.openai.azure.com', description: 'Azure hosted OpenAI' },
  { id: 'aws-bedrock', name: 'AWS Bedrock', type: 'aws-bedrock', baseUrl: 'https://bedrock-runtime.us-east-1.amazonaws.com', description: 'Claude, Llama, Titan' },
  { id: 'vertex-ai', name: 'Google Vertex AI', type: 'vertexai', baseUrl: 'https://us-central1-aiplatform.googleapis.com', description: 'Enterprise AI platform' },
  { id: 'ollama', name: 'Ollama', type: 'ollama', baseUrl: 'http://localhost:11434', description: 'Local LLM runtime' },
  { id: 'mistral', name: 'Mistral AI', type: 'openai', baseUrl: 'https://api.mistral.ai/v1', description: 'Mistral, Mixtral' },
  { id: 'cohere', name: 'Cohere', type: 'openai', baseUrl: 'https://api.cohere.ai/v1', description: 'Command, Embed' },
  { id: 'replicate', name: 'Replicate', type: 'openai', baseUrl: 'https://api.replicate.com/v1', description: 'Various open-source models' },
  { id: 'together', name: 'Together AI', type: 'openai', baseUrl: 'https://api.together.xyz/v1', description: 'Open-source LLMs' },
  { id: 'groq', name: 'Groq', type: 'openai', baseUrl: 'https://api.groq.com/openai/v1', description: 'Fast inference' },
  { id: 'perplexity', name: 'Perplexity', type: 'openai', baseUrl: 'https://api.perplexity.ai', description: 'AI search' },
  { id: 'deepseek', name: 'DeepSeek', type: 'openai', baseUrl: 'https://api.deepseek.com/v1', description: 'DeepSeek models' },
  { id: 'moonshot', name: 'Moonshot AI', type: 'openai', baseUrl: 'https://api.moonshot.cn/v1', description: 'Kimi models' },
  { id: 'zhipu', name: 'Zhipu AI', type: 'openai', baseUrl: 'https://open.bigmodel.cn/api/paas/v4', description: 'GLM models' },
  { id: 'baidu', name: 'Baidu ERNIE', type: 'openai', baseUrl: 'https://aip.baidubce.com/rpc/2.0/ai_custom/v1', description: 'ERNIE, Wenxin' },
  { id: 'alibaba', name: 'Alibaba Qwen', type: 'openai', baseUrl: 'https://dashscope.aliyuncs.com/api/v1', description: 'Qwen, Tongyi' },
  { id: 'minimax', name: 'MiniMax', type: 'openai', baseUrl: 'https://api.minimax.chat/v1', description: 'abab models' },
  { id: 'yi', name: '01.AI', type: 'openai', baseUrl: 'https://api.lingyiwanwu.com/v1', description: 'Yi models' },
  { id: 'baichuan', name: 'Baichuan', type: 'openai', baseUrl: 'https://api.baichuan-ai.com/v1', description: 'Baichuan models' },
  { id: 'sensetime', name: 'SenseTime', type: 'openai', baseUrl: 'https://api.sensenova.cn/v1', description: 'SenseChat' },
  { id: 'xunfei', name: 'iFlytek', type: 'openai', baseUrl: 'https://spark-api-open.xf-yun.com/v1', description: 'Spark models' },
  { id: 'huggingface', name: 'Hugging Face', type: 'openai', baseUrl: 'https://api-inference.huggingface.co/models', description: 'Inference API' },
  { id: 'novita', name: 'Novita AI', type: 'openai', baseUrl: 'https://api.novita.ai/v3/openai', description: 'LLM hosting' },
  { id: 'siliconflow', name: 'SiliconFlow', type: 'openai', baseUrl: 'https://api.siliconflow.cn/v1', description: 'Model hosting' },
  { id: 'fireworks', name: 'Fireworks AI', type: 'openai', baseUrl: 'https://api.fireworks.ai/inference/v1', description: 'Fast inference' },
  { id: 'anyscale', name: 'Anyscale', type: 'openai', baseUrl: 'https://api.endpoints.anyscale.com/v1', description: 'Ray-based hosting' },
  { id: 'lepton', name: 'Lepton AI', type: 'openai', baseUrl: 'https://api.lepton.ai/v1', description: 'AI deployment' },
  { id: 'modal', name: 'Modal', type: 'openai', baseUrl: 'https://api.modal.ai/v1', description: 'Serverless AI' },
  { id: 'runpod', name: 'RunPod', type: 'openai', baseUrl: 'https://api.runpod.ai/v1', description: 'GPU rental' },
  { id: 'vast', name: 'Vast AI', type: 'openai', baseUrl: 'https://api.vast.ai/v1', description: 'GPU marketplace' },
  { id: 'lambda', name: 'Lambda Labs', type: 'openai', baseUrl: 'https://api.lambdalabs.com/v1', description: 'GPU cloud' },
  { id: 'cerebras', name: 'Cerebras', type: 'openai', baseUrl: 'https://api.cerebras.ai/v1', description: 'Wafer-scale inference' },
  { id: 'sambanova', name: 'SambaNova', type: 'openai', baseUrl: 'https://api.sambanova.ai/v1', description: 'Enterprise AI' },
  { id: 'grok', name: 'xAI Grok', type: 'openai', baseUrl: 'https://api.x.ai/v1', description: 'Grok models' },
  { id: 'inflection', name: 'Inflection AI', type: 'openai', baseUrl: 'https://api.inflection.ai/v1', description: 'Pi assistant' },
  { id: 'character', name: 'Character.AI', type: 'openai', baseUrl: 'https://api.character.ai/v1', description: 'Character chat' },
  { id: 'stability', name: 'Stability AI', type: 'openai', baseUrl: 'https://api.stability.ai/v1', description: 'Stable Diffusion' },
  { id: 'midjourney', name: 'Midjourney', type: 'openai', baseUrl: 'https://api.midjourney.com/v1', description: 'Image generation' },
  { id: 'leonardo', name: 'Leonardo AI', type: 'openai', baseUrl: 'https://cloud.leonardo.ai/api/v1', description: 'Creative AI' },
  { id: 'ideogram', name: 'Ideogram', type: 'openai', baseUrl: 'https://api.ideogram.ai/v1', description: 'Text in images' },
  { id: 'flux', name: 'Flux', type: 'openai', baseUrl: 'https://api.flux.ai/v1', description: 'Image generation' },
  { id: 'jina', name: 'Jina AI', type: 'openai', baseUrl: 'https://api.jina.ai/v1', description: 'Embeddings, Rerank' },
  { id: 'voyage', name: 'Voyage AI', type: 'openai', baseUrl: 'https://api.voyageai.com/v1', description: 'Embeddings' },
  { id: 'pinecone', name: 'Pinecone', type: 'openai', baseUrl: 'https://api.pinecone.io', description: 'Vector database' },
  { id: 'weaviate', name: 'Weaviate', type: 'openai', baseUrl: 'https://api.weaviate.io/v1', description: 'Vector search' },
  { id: 'qdrant', name: 'Qdrant', type: 'openai', baseUrl: 'https://api.qdrant.io/v1', description: 'Vector DB' },
  { id: 'milvus', name: 'Milvus', type: 'openai', baseUrl: 'https://api.milvus.io/v1', description: 'Vector database' },
  { id: 'chroma', name: 'Chroma', type: 'openai', baseUrl: 'https://api.trychroma.com/v1', description: 'Embedding DB' },
  { id: 'elevenlabs', name: 'ElevenLabs', type: 'openai', baseUrl: 'https://api.elevenlabs.io/v1', description: 'Voice AI' },
  { id: 'assemblyai', name: 'AssemblyAI', type: 'openai', baseUrl: 'https://api.assemblyai.com/v2', description: 'Speech-to-text' },
  { id: 'deepgram', name: 'Deepgram', type: 'openai', baseUrl: 'https://api.deepgram.com/v1', description: 'Audio AI' },
  { id: 'whisper', name: 'OpenAI Whisper', type: 'openai', baseUrl: 'https://api.openai.com/v1', description: 'Speech recognition' },
  { id: 'tavily', name: 'Tavily', type: 'openai', baseUrl: 'https://api.tavily.com/v1', description: 'Web search API' },
  { id: 'serper', name: 'Serper', type: 'openai', baseUrl: 'https://api.serper.dev/v1', description: 'Google search' },
  { id: 'bing', name: 'Bing Search', type: 'openai', baseUrl: 'https://api.bing.microsoft.com/v7.0', description: 'Microsoft search' },
  { id: 'custom', name: 'Custom Provider', type: 'openai', baseUrl: '', description: 'Custom OpenAI-compatible API' }
]

export const getProviderById = (id: string): PredefinedProvider | undefined => {
  return predefinedProviders.find(p => p.id === id)
}
