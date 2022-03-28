export interface Integration {
  id: string
  api_token?: string
  created_at: string
  instance_name: string
  instance_url: string
  synced_at?: string
  syncing_error?: string
}
