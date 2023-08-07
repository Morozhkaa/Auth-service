package memory

import (
	"auth/internal/domain/models"
)

type MemoryStorage struct {
	storage map[string]models.User
}

func New() (*MemoryStorage, error) {
	storage := make(map[string]models.User)
	storage["Olenka"] = models.User{Login: "Olenka", Password: "Z55+a4XKM9tBClCi4/z9soOEK1/th6bWGveVqhgZTth3uoXAt+afxpy9m77Mo+y7LHJVKipxbOQL1u90V9oceaHaQATc9DH5UB8SEtYg/I6NKyrrnjQasdy7NBN++6834ZErQEsA6+9DmIr4ER3H2ecnQbXiRjBHQ5M2hzvTqc8=", Email: "Olya@mail.ru"} // password: dAmNmO!nAoBiZPi
	storage["Katya"] = models.User{Login: "Katya", Password: "b7t2PnSfjgF7/7Rr+gvOd5whra5HP7q9bV6AXp5sdRfQN0R4ashgfSr6hXi8KxkWQVf3ebmOAngocSc6Wo9HOX/I6OxIACEptozQ4eOwC0PR15ZO3w5SlOWMe6+wyjaJwOdOjhcPHQ1cP5DxxkWlIY+p/7XjcqHUNMzYdYQss8I=", Email: "Katya@mail.ru"}  // password: !HMnrsQ4VaGnJ-kK

	return &MemoryStorage{
		storage: storage,
	}, nil
}
