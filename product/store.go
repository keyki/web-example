package product

import (
    "errors"
    "github.com/jackc/pgx/v5/pgconn"
    "gorm.io/gorm"
    "log"
)

type Store struct {
    db *gorm.DB
}

func NewStore(db *gorm.DB) *Store {
    store := &Store{db}
    store.init()
    return store
}

func (s *Store) ListAll() ([]*Product, error) {
    log.Printf("Listing products")
    var products []*Product
    result := s.db.Find(&products)
    if result.Error != nil {
        return products, result.Error
    }
    return products, nil
}

func (s *Store) Create(product *Product) error {
    log.Printf("Creating product: %v", *product)
    result := s.db.Create(&product)
    if result.Error != nil {
        return result.Error
    }
    return nil
}

func (s *Store) FindByName(name string) (*Product, error) {
    log.Printf("Finding user by username: %v", name)
    var product Product
    result := s.db.Where("name = ?", name).First(&product)
    if result.Error != nil {
        return nil, result.Error
    }
    log.Printf("Found product: %v", product)
    return &product, nil
}

func (s *Store) init() {
    err := s.Create(&Product{
        Name:        "Pen",
        Description: "Blue pen",
        Price:       100,
        Currency:    HUF,
        Quantity:    50,
    })
    if err == nil {
        return
    }
    var pgErr *pgconn.PgError
    if errors.As(err, &pgErr) {
        if pgErr.Code == "23505" {
            log.Printf("Product already exists")
        } else {
            log.Println(err)
        }
    } else {
        log.Printf("Cannot create product: %v", err)
    }
}
