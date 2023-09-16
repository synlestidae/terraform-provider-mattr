variable "mattr_api_url" {
  description = "MATTR API URL"
  type        = string
  default     = null
}

variable "mattr_auth_url" {
  description = "MATTR Auth URL"
  type        = string
  default     = null
}

variable "mattr_client_id" {
  description = "MATTR Client ID"
  type        = string
  default     = null
}

variable "mattr_client_secret" {
  description = "MATTR Client Secret"
  type        = string
  default     = null
}

variable "mattr_auth_audience" {
  description = "MATTR Auth Audience"
  type        = string
  default     = null
}

variable "mattr_audience" {
  description = "MATTR Audience"
  type        = string
  default     = null
}

variable "mattr_access_token" {
  description = "Access token for API if you want to use a fixed access token"
  type        = string
  default     = null
}

variable "ngrok_url" {
  type = string
}
