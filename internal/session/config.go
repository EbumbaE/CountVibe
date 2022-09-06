package session

type Config struct{
    Paths map[string]string
    FormatsPages map[string]string
    JwtKey map[string][]byte
}

type ViewUserData struct{
    IsLogin bool
}

type Meal struct{
    TotalFats int64
    TotalCarbs int64
    TotalProts int64
    TotalCals int64
    Items []Item
}

type Item struct{
    ID int64
    Name string
    Amount int64
    Fat int64
    Carbs int64
    Prot int64
    Cals int64
}

type ViewDiaryData struct{
    IsLogin bool
    Date string
    TotalFats int64
    TotalCarbs int64
    TotalProts int64
    TotalCals int64
    
    Breakfast Meal
    Lunch Meal
    Dinner Meal
    Snacks Meal 
}
