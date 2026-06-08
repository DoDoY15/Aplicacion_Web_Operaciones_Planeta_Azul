export type UserRole = 'admin' | 'chefe_geral' | 'chefe_area' | 'supervisor' | 'membro'

export type ItemStatus = 'draft' | 'pending' | 'in_progress' | 'waiting' | 'done' | 'rejected'

export type ItemVisibility = 'private' | 'team' | 'public'

export type ItemPriority = 'low' | 'medium' | 'high' | 'urgent'

export interface Area {
  id: string
  name: string
  description?: string
  created_at: string
}

export interface User {
  id: string
  name: string
  email: string
  role: UserRole
  cargo?: string
  area_id?: string
  is_active: boolean
  created_at: string
  area?: Area
}

export interface Item {
  id: string
  parent_id?: string
  title: string
  description?: string
  created_by: string
  assigned_to?: string
  area_id: string
  status: ItemStatus
  visibility: ItemVisibility
  requires_approval: boolean
  priority: ItemPriority
  deadline?: string
  completed_at?: string
  created_at: string
  updated_at: string
  deleted_at?: string
  creator?: User
  assignee?: User
  children?: Item[]
  comments?: Comment[]
}

export interface Comment {
  id: string
  item_id: string
  user_id: string
  content: string
  created_at: string
  author?: User
}

export interface Interaction {
  id: string
  item_id: string
  opened_by: string
  addressed_to: string
  message: string
  response?: string
  status: 'open' | 'resolved'
  resolved_at?: string
  created_at: string
}

export interface Notification {
  id: string
  user_id: string
  type: string
  ref_id: string
  ref_type: string
  message: string
  read: boolean
  created_at: string
}

export interface LoginResponse {
  access_token: string
  refresh_token: string
  user: User
}

// Display helpers
export const STATUS_LABELS: Record<ItemStatus, string> = {
  draft:       'Borrador',
  pending:     'Pendiente',
  in_progress: 'En Progreso',
  waiting:     'En Espera',
  done:        'Completado',
  rejected:    'Rechazado',
}

export const PRIORITY_LABELS: Record<ItemPriority, string> = {
  low:    'Baja',
  medium: 'Media',
  high:   'Alta',
  urgent: 'Urgente',
}

export const ROLE_LABELS: Record<UserRole, string> = {
  admin:       'Administrador',
  chefe_geral: 'Jefe General',
  chefe_area:  'Jefe de Área',
  supervisor:  'Supervisor',
  membro:      'Miembro',
}
