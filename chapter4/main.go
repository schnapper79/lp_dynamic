// Exhaustive search

package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

const num_items = 40 // A reasonable value for exhaustive search.

const min_value = 1
const max_value = 10
const min_weight = 4
const max_weight = 10

var allowed_weight int

type Item struct {
	id, blocked_by int
	block_list     []int // Other items that this one blocks.
	value, weight  int
	is_selected    bool
}

// Make some random items.
func make_items(num_items, min_value, max_value, min_weight, max_weight int) []Item {
	// Initialize a pseudorandom number generator.
	//random := rand.New(rand.NewSource(time.Now().UnixNano())) // Initialize with a changing seed
	random := rand.New(rand.NewSource(1337)) // Initialize with a fixed seed

	items := make([]Item, num_items)
	for i := 0; i < num_items; i++ {
		items[i] = Item{
			i, -1, nil,
			random.Intn(max_value-min_value+1) + min_value,
			random.Intn(max_weight-min_weight+1) + min_weight,
			false}
	}
	return items
}

// Return a copy of the items slice.
func copy_items(items []Item) []Item {
	new_items := make([]Item, len(items))
	copy(new_items, items)
	return new_items
}

// Return the total value of the items.
// If add_all is false, only add up the selected items.
func sum_values(items []Item, add_all bool) int {
	total := 0
	for i := 0; i < len(items); i++ {
		if add_all || items[i].is_selected {
			total += items[i].value
		}
	}
	return total
}

// Return the total weight of the items.
// If add_all is false, only add up the selected items.
func sum_weights(items []Item, add_all bool) int {
	total := 0
	for i := 0; i < len(items); i++ {
		if add_all || items[i].is_selected {
			total += items[i].weight
		}
	}
	return total
}

// Return the value of this solution.
// If the solution is too heavy, return -1 so we prefer an empty solution.
func solution_value(items []Item, allowed_weight int) int {
	// If the solution's total weight > allowed_weight,
	// return -1 so we won't use this solution.
	if sum_weights(items, false) > allowed_weight {
		return -1
	}

	// Return the sum of the selected values.
	return sum_values(items, false)
}

// Print the selected items.
func print_selected(items []Item) {
	num_printed := 0
	for i, item := range items {
		if item.is_selected {
			fmt.Printf("%d(%d, %d) ", i, item.value, item.weight)
		}
		num_printed += 1
		if num_printed > 100 {
			fmt.Println("...")
			return
		}
	}
	fmt.Println()
}

func run_algorithm(alg func([]Item, int) ([]Item, int, int), items []Item, allowed_weight int) {
	// Copy the items so the run isn't influenced by a previous run.
	test_items := copy_items(items)

	start := time.Now()

	// Run the algorithm.
	solution, total_value, function_calls := alg(test_items, allowed_weight)

	elapsed := time.Since(start)

	fmt.Printf("Elapsed: %f\n", elapsed.Seconds())
	print_selected(solution)
	fmt.Printf("Value: %d, Weight: %d, Calls: %d\n",
		total_value, sum_weights(solution, false), function_calls)
	fmt.Println()
}

func exhaustive_search(items []Item, allowed_weight int) ([]Item, int, int) {
	return do_exhaustive_search(items, allowed_weight, 0)
}

func do_exhaustive_search(items []Item, allowed_weight, next_index int) ([]Item, int, int) {
	if next_index >= len(items) {
		return copy_items(items), solution_value(items, allowed_weight), 1
	}
	//try to add item
	items[next_index].is_selected = true
	best_items, best_value, function_calls := do_exhaustive_search(items, allowed_weight, next_index+1)
	//try to remove item
	items[next_index].is_selected = false
	other_items, other_value, other_calls := do_exhaustive_search(items, allowed_weight, next_index+1)
	function_calls += other_calls
	if other_value > best_value {
		best_items = other_items
		best_value = other_value
	}
	return best_items, best_value, function_calls + 1
}

func branch_and_bound(items []Item, allowed_weight int) ([]Item, int, int) {
	best_value := 0
	current_value := 0
	current_weight := 0
	remaing_value := 0
	for _, item := range items {
		remaing_value += item.value
	}

	return do_branch_and_bound(items, allowed_weight, 0, best_value, current_value, current_weight, remaing_value)
}

func do_branch_and_bound(items []Item, allowed_weight, next_index, best_value, current_value, current_weight, remaing_value int) ([]Item, int, int) {
	if next_index >= len(items) {
		copied_Items := copy_items(items)
		return copied_Items, current_value, 1
	}

	if current_value+remaing_value <= best_value {
		return nil, current_value, 1
	}

	var sol_items1 []Item
	var sol_value1 int
	var sol_calls1 int

	if current_weight+items[next_index].weight <= allowed_weight {
		items[next_index].is_selected = true
		sol_items1, sol_value1, sol_calls1 = do_branch_and_bound(items, allowed_weight, next_index+1, best_value, current_value+items[next_index].value, current_weight+items[next_index].weight, remaing_value-items[next_index].value)
		if sol_value1 > best_value {
			best_value = sol_value1
		}
	} else {
		sol_items1, sol_value1, sol_calls1 = nil, 0, 1
	}

	var sol_items2 []Item
	var sol_value2 int
	var sol_calls2 int

	items[next_index].is_selected = false
	sol_items2, sol_value2, sol_calls2 = do_branch_and_bound(items, allowed_weight, next_index+1, best_value, current_value, current_weight, remaing_value-items[next_index].value)

	sol_calls1 += sol_calls2
	if sol_value1 > sol_value2 {
		return sol_items1, sol_value1, sol_calls1 + 1
	} else {
		return sol_items2, sol_value2, sol_calls1 + 1
	}

}

func rods_technique(items []Item, allowed_weight int) ([]Item, int, int) {
	best_value := 0
	current_value := 0
	current_weight := 0
	remaing_value := 0
	for _, item := range items {
		remaing_value += item.value
	}

	make_block_lists(items)

	return do_rods_technique(items, allowed_weight, 0, best_value, current_value, current_weight, remaing_value)
}

func do_rods_technique(items []Item, allowed_weight, next_index, best_value, current_value, current_weight, remaing_value int) ([]Item, int, int) {
	if next_index >= len(items) {
		copied_Items := copy_items(items)
		return copied_Items, current_value, 1
	}

	if current_value+remaing_value <= best_value {
		return nil, current_value, 1
	}

	var sol_items1 []Item
	sol_items1 = nil
	sol_value1 := 0
	sol_calls1 := 0

	if items[next_index].blocked_by == -1 {
		if current_weight+items[next_index].weight <= allowed_weight {
			items[next_index].is_selected = true
			sol_items1, sol_value1, sol_calls1 = do_rods_technique(items, allowed_weight, next_index+1, best_value, current_value+items[next_index].value, current_weight+items[next_index].weight, remaing_value-items[next_index].value)
			if sol_value1 > best_value {
				best_value = sol_value1
			}
		}
	}
	var sol_items2 []Item
	var sol_value2 int
	var sol_calls2 int

	items[next_index].is_selected = false
	block_items(items[next_index], items)
	sol_items2, sol_value2, sol_calls2 = do_rods_technique(items, allowed_weight, next_index+1, best_value, current_value, current_weight, remaing_value-items[next_index].value)
	unblock_items(items[next_index], items)
	sol_calls1 += sol_calls2
	if sol_value1 > sol_value2 {
		return sol_items1, sol_value1, sol_calls1 + 1
	} else {

		return sol_items2, sol_value2, sol_calls1 + 1
	}
}

func rods_technique_sorted(items []Item, allowed_weight int) ([]Item, int, int) {
	best_value := 0
	current_value := 0
	current_weight := 0
	remaing_value := 0
	for _, item := range items {
		remaing_value += item.value
	}

	make_block_lists(items)
	// Sort so items with longer blocked lists come first.
	sort.Slice(items, func(i, j int) bool {
		return len(items[i].block_list) > len(items[j].block_list)
	})

	// Reset the items' IDs.
	for i := range items {
		items[i].id = i
	}

	// Rebuild the blocked lists with the new indices.
	make_block_lists(items)

	return do_rods_technique(items, allowed_weight, 0, best_value, current_value, current_weight, remaing_value)
}

func make_block_lists(items []Item) {
	for i, item := range items {
		items[i].block_list = make([]int, 0)
		for j, other_item := range items {
			if i == j {
				continue
			}
			if item.weight <= other_item.weight && item.value >= other_item.value {
				items[i].block_list = append(items[i].block_list, items[j].id)
			}
		}
	}
}

func block_items(source Item, items []Item) {
	for _, blocked_by := range source.block_list {
		if items[blocked_by].blocked_by == -1 {
			items[blocked_by].blocked_by = source.id
		}
	}
}

func unblock_items(source Item, items []Item) {
	for _, blocked_by := range source.block_list {
		if items[blocked_by].blocked_by == source.id {
			items[blocked_by].blocked_by = -1
		}
	}
}

func getSliceOfSlices(r, c int) [][]int {
	slice := make([][]int, r)
	for i := range slice {
		slice[i] = make([]int, c)
	}
	return slice
}

// Use dynamic programming to find a solution.
// Return the best assignment, value of that assignment,
// and the number of function calls we made.
func dynamic_programming(items []Item, allowed_weight int) ([]Item, int, int) {
	solution_value_array := getSliceOfSlices(len(items), allowed_weight+1)
	prev_weight_array := getSliceOfSlices(len(items), allowed_weight+1)

	//initialize first row
	for i := 0; i < allowed_weight+1; i++ {
		if items[0].weight <= i {
			solution_value_array[0][i] = items[0].value
			prev_weight_array[0][i] = -1
		} else {
			solution_value_array[0][i] = 0
			prev_weight_array[0][i] = i
		}
	}

	for i := 1; i < len(items); i++ {
		for j := 0; j < allowed_weight+1; j++ {
			//Calculate the value we get if item i is not in the solution.
			value_without_item := solution_value_array[i-1][j]
			//Calculate the value we get if item i is in the solution.
			value_with_item := 0
			if items[i].weight <= j {
				value_with_item = solution_value_array[i-1][j-items[i].weight] + items[i].value
			}
			//Choose the better of the two values.
			if value_with_item > value_without_item {
				solution_value_array[i][j] = value_with_item
				prev_weight_array[i][j] = j - items[i].weight
			} else {
				solution_value_array[i][j] = value_without_item
				prev_weight_array[i][j] = j
			}
		}
	}
	for i := 0; i < len(items); i++ {
		items[i].is_selected = false
	}
	//Find the items in the solution.
	i := len(items) - 1
	j := allowed_weight
	for i >= 0 {
		if prev_weight_array[i][j] != j {
			items[i].is_selected = true
			j = prev_weight_array[i][j]
		}
		i--
	}
	return items, sum_values(items, false), 1
}

func main() {
	items := make_items(num_items, min_value, max_value, min_weight, max_weight)
	allowed_weight = sum_weights(items, true) / 2

	// Display basic parameters.
	fmt.Println("*** Parameters ***")
	fmt.Printf("# items: %d\n", num_items)
	fmt.Printf("Total value: %d\n", sum_values(items, true))
	fmt.Printf("Total weight: %d\n", sum_weights(items, true))
	fmt.Printf("Allowed weight: %d\n", allowed_weight)
	fmt.Println()

	// Exhaustive search
	if num_items > 25 { // Only run exhaustive search if num_items <= 25.
		fmt.Println("Too many items for exhaustive search\n")
	} else {
		fmt.Println("*** Exhaustive Search ***")
		run_algorithm(exhaustive_search, items, allowed_weight)
	}

	// branch_and_bound search
	if num_items > 45 { // Only run branch_and_bound search if num_items <= 25.
		fmt.Println("Too many items for branch_and_bound search\n")
	} else {
		fmt.Println("*** branch_and_bound ***")
		run_algorithm(branch_and_bound, items, allowed_weight)
	}
	// Rod's technique
	if num_items > 85 { // Only use Rod's technique if num_items <= 85.
		fmt.Println("Too many items for Rod's technique\n")
	} else {
		fmt.Println("*** Rod's technique ***")
		run_algorithm(rods_technique, items, allowed_weight)
	}
	// Rod's sorted technique
	if num_items > 350 { // Only use Rod's technique if num_items <= 85.
		fmt.Println("Too many items for Rod's sorted  technique\n")
	} else {
		fmt.Println("*** Rod's sorted technique ***")
		run_algorithm(rods_technique_sorted, items, allowed_weight)
	}
	// Dynamic programming
	fmt.Println("*** Dynamic programming ***")
	run_algorithm(dynamic_programming, items, allowed_weight)
}
