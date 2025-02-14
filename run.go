// Copyright (c) 2025 Michael D Henderson. All rights reserved.

package main

//func run(path string, debug bool) error {
//	e, err := ec.LoadGame(path)
//	if err != nil {
//		return err
//	}
//
//	// load all the files
//	for _, name := range []string{"orders.txt"} {
//		input, err := os.ReadFile(name)
//		if err != nil {
//			return err
//		}
//		lexemes, err := orders.Scan(input)
//		if err != nil {
//			return err
//		}
//		ods := orders.Parse(lexemes)
//		if debug {
//			for _, od := range ods {
//				fmt.Println(od)
//			}
//		}
//		err = e.AddOrders(adapters.OrdersToEngineOrders(ods))
//		if err != nil {
//			log.Printf("%s: %v\n", name, err)
//		}
//	}
//
//	err = e.Process()
//	if err != nil {
//		return err
//	}
//
//	return e.SaveGame(path)
//}
