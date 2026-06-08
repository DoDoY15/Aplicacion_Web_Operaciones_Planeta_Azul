package repository

import (
	"errors"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/planeta-azul/backend/internal/models"
	"golang.org/x/crypto/bcrypt"
)

// MemStore is a thread-safe in-memory store for running without a database.
// Replace with real DB repository when PostgreSQL is available.
type MemStore struct {
	mu            sync.RWMutex
	users         map[uuid.UUID]*models.User
	items         map[uuid.UUID]*models.Item
	comments      map[uuid.UUID]*models.Comment
	interactions  map[uuid.UUID]*models.Interaction
	approvals     map[uuid.UUID]*models.Approval
	notifications map[uuid.UUID]*models.Notification
	areas         map[uuid.UUID]*models.Area
}

func NewMemStore() *MemStore {
	store := &MemStore{
		users:         make(map[uuid.UUID]*models.User),
		items:         make(map[uuid.UUID]*models.Item),
		comments:      make(map[uuid.UUID]*models.Comment),
		interactions:  make(map[uuid.UUID]*models.Interaction),
		approvals:     make(map[uuid.UUID]*models.Approval),
		notifications: make(map[uuid.UUID]*models.Notification),
		areas:         make(map[uuid.UUID]*models.Area),
	}
	store.seed()
	return store
}

func (s *MemStore) seed() {
	// Areas
	areaA := &models.Area{ID: uuid.MustParse("11111111-0000-0000-0000-000000000001"), Name: "Producción", Description: "Línea de producción principal", CreatedAt: time.Now()}
	areaB := &models.Area{ID: uuid.MustParse("11111111-0000-0000-0000-000000000002"), Name: "Mantenimiento", Description: "Mantenimiento industrial", CreatedAt: time.Now()}
	s.areas[areaA.ID] = areaA
	s.areas[areaB.ID] = areaB

	hashAdmin, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	hashChefe, _ := bcrypt.GenerateFromPassword([]byte("chefe123"), bcrypt.DefaultCost)
	hashSup, _ := bcrypt.GenerateFromPassword([]byte("sup123"), bcrypt.DefaultCost)
	hashMembro, _ := bcrypt.GenerateFromPassword([]byte("membro123"), bcrypt.DefaultCost)

	adminID := uuid.MustParse("22222222-0000-0000-0000-000000000001")
	chefeID := uuid.MustParse("22222222-0000-0000-0000-000000000002")
	supID := uuid.MustParse("22222222-0000-0000-0000-000000000003")
	membroID := uuid.MustParse("22222222-0000-0000-0000-000000000004")

	s.users[adminID] = &models.User{
		ID: adminID, Name: "Admin Sistema", Email: "admin@planetaazul.com",
		PasswordHash: string(hashAdmin), Role: models.RoleAdmin, IsActive: true, CreatedAt: time.Now(),
	}
	s.users[chefeID] = &models.User{
		ID: chefeID, Name: "Carlos Mendoza", Email: "chefe@planetaazul.com",
		PasswordHash: string(hashChefe), Role: models.RoleChefArea, Cargo: "Jefe de Producción",
		AreaID: &areaA.ID, IsActive: true, CreatedAt: time.Now(),
	}
	s.users[supID] = &models.User{
		ID: supID, Name: "Ana Ferreira", Email: "sup@planetaazul.com",
		PasswordHash: string(hashSup), Role: models.RoleSupervisor, Cargo: "Supervisora de Turno",
		AreaID: &areaA.ID, IsActive: true, CreatedAt: time.Now(),
	}
	s.users[membroID] = &models.User{
		ID: membroID, Name: "Juan Silva", Email: "membro@planetaazul.com",
		PasswordHash: string(hashMembro), Role: models.RoleMembro, Cargo: "Operador",
		AreaID: &areaA.ID, IsActive: true, CreatedAt: time.Now(),
	}

	// Seed items
	deadline := time.Now().Add(72 * time.Hour)
	rootItem := &models.Item{
		ID: uuid.MustParse("33333333-0000-0000-0000-000000000001"),
		Title: "Mantenimiento Preventivo — Línea A", Description: "Revisión completa de los equipos de la línea A antes de la parada de julio.",
		CreatedBy: chefeID, AssignedTo: &supID, AreaID: areaA.ID,
		Status: models.StatusInProgress, Visibility: models.VisibilityTeam,
		Priority: models.PriorityHigh, RequiresApproval: false,
		Deadline: &deadline, CreatedAt: time.Now().Add(-24 * time.Hour), UpdatedAt: time.Now(),
	}
	child1ID := uuid.MustParse("33333333-0000-0000-0000-000000000002")
	child1 := &models.Item{
		ID: child1ID, ParentID: &rootItem.ID,
		Title: "Revisar la lubricación de los rodamientos", CreatedBy: supID,
		AssignedTo: &membroID, AreaID: areaA.ID,
		Status: models.StatusDone, Visibility: models.VisibilityTeam,
		Priority: models.PriorityMedium, CreatedAt: time.Now().Add(-20 * time.Hour), UpdatedAt: time.Now(),
	}
	child2 := &models.Item{
		ID: uuid.MustParse("33333333-0000-0000-0000-000000000003"), ParentID: &rootItem.ID,
		Title: "Sustituir los filtros del sistema hidráulico", CreatedBy: supID,
		AssignedTo: &membroID, AreaID: areaA.ID,
		Status: models.StatusInProgress, Visibility: models.VisibilityTeam,
		Priority: models.PriorityHigh, Deadline: &deadline,
		CreatedAt: time.Now().Add(-18 * time.Hour), UpdatedAt: time.Now(),
	}
	item2 := &models.Item{
		ID: uuid.MustParse("33333333-0000-0000-0000-000000000004"),
		Title: "Solicitud de EPP — Turno de Tarde", Description: "Solicitar reposición de cascos y guantes para el turno de tarde.",
		CreatedBy: membroID, AssignedTo: &supID, AreaID: areaA.ID,
		Status: models.StatusPending, Visibility: models.VisibilityTeam,
		Priority: models.PriorityMedium, RequiresApproval: true,
		CreatedAt: time.Now().Add(-2 * time.Hour), UpdatedAt: time.Now(),
	}
	item3 := &models.Item{
		ID: uuid.MustParse("33333333-0000-0000-0000-000000000005"),
		Title: "Informe de Producción — Semana 23", Description: "Compilar los datos de producción de la semana 23 para la presentación al jefe general.",
		CreatedBy: chefeID, AreaID: areaA.ID,
		Status: models.StatusDraft, Visibility: models.VisibilityPrivate,
		Priority: models.PriorityLow,
		CreatedAt: time.Now().Add(-5 * time.Hour), UpdatedAt: time.Now(),
	}

	s.items[rootItem.ID] = rootItem
	s.items[child1ID] = child1
	s.items[child2.ID] = child2
	s.items[item2.ID] = item2
	s.items[item3.ID] = item3

	// Comments
	comment := &models.Comment{
		ID: uuid.New(), ItemID: rootItem.ID, UserID: supID,
		Content: "Inicié la inspección visual. Los rodamientos de la estación 3 presentan un desgaste superior al normal.",
		CreatedAt: time.Now().Add(-12 * time.Hour),
	}
	s.comments[comment.ID] = comment
}

// --- User methods ---

func (s *MemStore) GetUserByEmail(email string) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, u := range s.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, errors.New("usuario no encontrado")
}

func (s *MemStore) GetUserByID(id uuid.UUID) (*models.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	u, ok := s.users[id]
	if !ok {
		return nil, errors.New("usuario no encontrado")
	}
	return u, nil
}

func (s *MemStore) ListUsers() []*models.User {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*models.User, 0, len(s.users))
	for _, u := range s.users {
		out = append(out, u)
	}
	return out
}

// --- Item methods ---

func (s *MemStore) GetItemByID(id uuid.UUID) (*models.Item, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	item, ok := s.items[id]
	if !ok || item.DeletedAt != nil {
		return nil, errors.New("ítem no encontrado")
	}
	return item, nil
}

func (s *MemStore) ListItems(areaID *uuid.UUID, userID uuid.UUID, role models.UserRole) []*models.Item {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var out []*models.Item
	for _, item := range s.items {
		if item.DeletedAt != nil {
			continue
		}
		// Admins and chefe_geral see everything
		if role == models.RoleAdmin || role == models.RoleChefGeral {
			out = append(out, item)
			continue
		}
		// Others: filter by area
		if areaID != nil && item.AreaID == *areaID {
			// Private items only visible to creator or assignee
			if item.Visibility == models.VisibilityPrivate {
				if item.CreatedBy == userID || (item.AssignedTo != nil && *item.AssignedTo == userID) {
					out = append(out, item)
				}
			} else {
				out = append(out, item)
			}
		}
	}
	return out
}

func (s *MemStore) CreateItem(item *models.Item) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	item.ID = uuid.New()
	item.CreatedAt = time.Now()
	item.UpdatedAt = time.Now()
	s.items[item.ID] = item
	return nil
}

func (s *MemStore) UpdateItem(id uuid.UUID, updates map[string]interface{}) (*models.Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[id]
	if !ok || item.DeletedAt != nil {
		return nil, errors.New("ítem no encontrado")
	}
	if v, ok := updates["title"]; ok {
		item.Title = v.(string)
	}
	if v, ok := updates["description"]; ok {
		item.Description = v.(string)
	}
	if v, ok := updates["status"]; ok {
		item.Status = v.(models.ItemStatus)
	}
	if v, ok := updates["priority"]; ok {
		item.Priority = v.(models.ItemPriority)
	}
	if v, ok := updates["assigned_to"]; ok {
		id := v.(uuid.UUID)
		item.AssignedTo = &id
	}
	item.UpdatedAt = time.Now()
	return item, nil
}

func (s *MemStore) SoftDeleteItem(id uuid.UUID, deletedBy uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	item, ok := s.items[id]
	if !ok {
		return errors.New("ítem no encontrado")
	}
	now := time.Now()
	item.DeletedAt = &now
	item.DeletedBy = &deletedBy
	return nil
}

// --- Comment methods ---

func (s *MemStore) ListCommentsByItem(itemID uuid.UUID) []*models.Comment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*models.Comment
	for _, c := range s.comments {
		if c.ItemID == itemID {
			out = append(out, c)
		}
	}
	return out
}

func (s *MemStore) CreateComment(c *models.Comment) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	c.ID = uuid.New()
	c.CreatedAt = time.Now()
	s.comments[c.ID] = c
	return nil
}

// --- Interaction methods ---

func (s *MemStore) GetInteractionsByItem(itemID uuid.UUID) []*models.Interaction {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*models.Interaction
	for _, i := range s.interactions {
		if i.ItemID == itemID {
			out = append(out, i)
		}
	}
	return out
}

func (s *MemStore) CreateInteraction(i *models.Interaction) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	i.ID = uuid.New()
	i.CreatedAt = time.Now()
	i.Status = models.InteractionOpen
	s.interactions[i.ID] = i
	return nil
}

func (s *MemStore) ResolveInteraction(id uuid.UUID, response string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	i, ok := s.interactions[id]
	if !ok {
		return errors.New("interacción no encontrada")
	}
	now := time.Now()
	i.Response = response
	i.Status = models.InteractionResolved
	i.ResolvedAt = &now
	return nil
}

// --- Notification methods ---

func (s *MemStore) GetNotificationsByUser(userID uuid.UUID) []*models.Notification {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var out []*models.Notification
	for _, n := range s.notifications {
		if n.UserID == userID {
			out = append(out, n)
		}
	}
	return out
}

func (s *MemStore) CreateNotification(n *models.Notification) {
	s.mu.Lock()
	defer s.mu.Unlock()
	n.ID = uuid.New()
	n.CreatedAt = time.Now()
	s.notifications[n.ID] = n
}

func (s *MemStore) MarkNotificationRead(id uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if n, ok := s.notifications[id]; ok {
		n.Read = true
	}
}

// --- Area methods ---

func (s *MemStore) ListAreas() []*models.Area {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*models.Area, 0, len(s.areas))
	for _, a := range s.areas {
		out = append(out, a)
	}
	return out
}
