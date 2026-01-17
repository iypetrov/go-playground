package main

import "fmt"

type account struct {
	balance float32
}

type Customer struct {
	ID      string
	Balance float64
}

type Store struct {
	m map[string]*Customer
}

// after Go 1.22 this works as expected
func (s *Store) storeCustomersV1(customers []Customer) {
	for _, customer := range customers {
		fmt.Printf("%p\n", &customer)
		s.m[customer.ID] = &customer
	}
}

func main() {
	accounts := []account{
		{balance: 100.},
		{balance: 200.},
		{balance: 300.},
	}
	// for _, a := range accounts {
	// 	a.balance += 1000
	// }
	for i := range accounts {
		accounts[i].balance += 1000
	}
	fmt.Println(accounts)

	// s1 := []int{0, 1, 2}
	// for range s1 {
	// 	s1 = append(s1, 10)
	// }
	//for i := 0; i < len(s1); i++ {
	//	s1 = append(s1, 10)
	//}
	//fmt.Println(s1)
	// a := [3]int{0, 1, 2}
	// for i, v := range a {
	// 	a[2] = 10
	// 	if i == 2 {
	// 		fmt.Println(v)
	// 	}
	// }
	// for i, v := range &a {
	// 	a[2] = 10
	// 	if i == 2 {
	// 		fmt.Println(v)
	// 	}
	// }
	// fmt.Println(a)
	// s := Store{
	// 	m: make(map[string]*Customer),
	// }
	// s.storeCustomersV1([]Customer{
	// 	{ID: "1", Balance: 10},
	// 	{ID: "2", Balance: -10},
	// 	{ID: "3", Balance: 0},
	// })
	// for k, v := range s.m {
	// 	fmt.Printf("key=%s, value=%#v\n", k, v)
	// }
}
