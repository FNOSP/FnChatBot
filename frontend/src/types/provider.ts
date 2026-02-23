export type ProviderType =
  | 'openai'
  | 'openai-response'
  | 'anthropic'
  | 'gemini'
  | 'azure-openai'
  | 'vertexai'
  | 'mistral'
  | 'aws-bedrock'
  | 'vertex-anthropic'
  | 'new-api'
  | 'gateway'
  | 'ollama'

export type ModelCapability =
  | 'text'
  | 'vision'
  | 'embedding'
  | 'reasoning'
  | 'function_calling'
  | 'web_search'
  | 'rerank'
  | 'image_generation'

export type EndpointType =
  | 'openai'
  | 'openai-response'
  | 'anthropic'
  | 'gemini'
  | 'image-generation'
  | 'jina-rerank'

export interface ProviderApiOptions {
  is_not_support_array_content?: boolean
  is_not_support_stream_options?: boolean
  is_support_developer_role?: boolean
  is_support_service_tier?: boolean
  is_not_support_enable_thinking?: boolean
  is_not_support_api_version?: boolean
  is_not_support_verbosity?: boolean
}

export interface Model {
  id: number
  provider_id: number
  provider?: Provider
  model_id: string
  name: string
  group?: string
  description?: string
  owned_by?: string
  capabilities?: ModelCapability[]
  supported_endpoint_types?: EndpointType[]
  endpoint_type?: EndpointType
  max_tokens?: number
  input_price?: number
  output_price?: number
  supported_text_delta?: boolean
  enabled?: boolean
  is_default?: boolean
  created_at?: string
  updated_at?: string
}

export interface Provider {
  id: number
  provider_id: string
  name: string
  type: ProviderType
  base_url: string
  api_key?: string
  enabled: boolean
  is_system: boolean
  api_options?: ProviderApiOptions
  models?: Model[]
  created_at?: string
  updated_at?: string
}

export interface ModelInfo {
  id: string
  name: string
  owned_by?: string
  description?: string
  capabilities?: ModelCapability[]
}
