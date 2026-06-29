# LeetCode Solutions — Go

39 coding problems, each in its own folder as a standalone `package main`.

**Run any problem:**
```bash
go run leetcode/<folder>/<file>.go
# e.g.
go run leetcode/single_number/single_number.go
```

---

## Problems

| # | Problem | Difficulty | File | Approach |
|---|---------|------------|------|----------|
| 118 | Pascal's Triangle | Easy | [pascals_triangle/pascals_triangle.go](pascals_triangle/pascals_triangle.go) | DP — build each row from previous |
| 64 | Minimum Path Sum | Medium | [minimum_path_sum/minimum_path_sum.go](minimum_path_sum/minimum_path_sum.go) | DP — in-place grid update |
| 1710 | Maximum Units on a Truck | Easy | [maximum_units_on_a_truck/maximum_units_on_a_truck.go](maximum_units_on_a_truck/maximum_units_on_a_truck.go) | Greedy — sort by units desc |
| 205 | Isomorphic Strings | Easy | [isomorphic_strings/isomorphic_strings.go](isomorphic_strings/isomorphic_strings.go) | Two-way character mapping |
| 791 | Custom Sort String | Medium | [custom_sort_string/custom_sort_string.go](custom_sort_string/custom_sort_string.go) | Sort with rank map |
| 103 | Binary Tree Zigzag Level Order Traversal | Medium | [binary_tree_zigzag_level_order_traversal/binary_tree_zigzag_level_order_traversal.go](binary_tree_zigzag_level_order_traversal/binary_tree_zigzag_level_order_traversal.go) | BFS with alternating fill direction |
| 126 | Word Ladder II | Hard | [word_ladder_ii/word_ladder_ii.go](word_ladder_ii/word_ladder_ii.go) | BFS (by layer) + DFS path reconstruction |
| 199 | Binary Tree Right Side View | Medium | [binary_tree_right_side_view/binary_tree_right_side_view.go](binary_tree_right_side_view/binary_tree_right_side_view.go) | BFS — last node per level |
| 426 | Convert BST to Sorted Doubly Linked List | Medium | [convert_bst_to_sorted_doubly_linked_list/convert_bst_to_sorted_doubly_linked_list.go](convert_bst_to_sorted_doubly_linked_list/convert_bst_to_sorted_doubly_linked_list.go) | In-order traversal with prev pointer |
| 166 | Fraction to Recurring Decimal | Medium | [fraction_to_recurring_decimal/fraction_to_recurring_decimal.go](fraction_to_recurring_decimal/fraction_to_recurring_decimal.go) | Long division with remainder tracking |
| 58 | Length of Last Word | Easy | [length_of_last_word/length_of_last_word.go](length_of_last_word/length_of_last_word.go) | Scan from end, skip spaces |
| 136 | Single Number | Easy | [single_number/single_number.go](single_number/single_number.go) | XOR all elements |
| 815 | Bus Routes | Hard | [bus_routes/bus_routes.go](bus_routes/bus_routes.go) | BFS on routes (not stops) |
| 291 | Word Pattern II | Medium | [word_pattern_ii/word_pattern_ii.go](word_pattern_ii/word_pattern_ii.go) | Backtracking with two-way binding |
| 498 | Diagonal Traverse | Medium | [diagonal_traverse/diagonal_traverse.go](diagonal_traverse/diagonal_traverse.go) | Diagonal index math with direction alternation |
| 443 | String Compression | Medium | [string_compression/string_compression.go](string_compression/string_compression.go) | Two-pointer in-place run-length encoding |
| 8 | String to Integer (atoi) | Medium | [string_to_integer_atoi/string_to_integer_atoi.go](string_to_integer_atoi/string_to_integer_atoi.go) | Iterative parse with overflow clamp |
| 1366 | Rank Teams by Votes | Medium | [rank_teams_by_votes/rank_teams_by_votes.go](rank_teams_by_votes/rank_teams_by_votes.go) | Count votes per position, custom sort |
| 432 | All O'one Data Structure | Hard | [all_oone_data_structure/all_oone_data_structure.go](all_oone_data_structure/all_oone_data_structure.go) | Doubly linked list of count-buckets + hashmap |
| 11 | Container With Most Water | Medium | [container_with_most_water/container_with_most_water.go](container_with_most_water/container_with_most_water.go) | Two pointers, shrink shorter side |
| 378 | Kth Smallest Element in a Sorted Matrix | Medium | [kth_smallest_element_in_sorted_matrix/kth_smallest_element_in_sorted_matrix.go](kth_smallest_element_in_sorted_matrix/kth_smallest_element_in_sorted_matrix.go) | Binary search on value + count ≤ mid |
| 9 | Palindrome Number | Easy | [palindrome_number/palindrome_number.go](palindrome_number/palindrome_number.go) | Reverse second half, compare halves |
| 234 | Palindrome Linked List | Easy | [palindrome_linked_list/palindrome_linked_list.go](palindrome_linked_list/palindrome_linked_list.go) | Find mid, reverse second half, compare |
| 1455 | Check If a Word Occurs As a Prefix of Any Word in a Sentence | Easy | [check_word_prefix_in_sentence/check_word_prefix_in_sentence.go](check_word_prefix_in_sentence/check_word_prefix_in_sentence.go) | Split + HasPrefix |
| 3355 | Zero Array Transformation I | Medium | [zero_array_transformation_i/zero_array_transformation_i.go](zero_array_transformation_i/zero_array_transformation_i.go) | Difference array + prefix sum |
| 329 | Longest Increasing Path in a Matrix | Hard | [longest_increasing_path_in_matrix/longest_increasing_path_in_matrix.go](longest_increasing_path_in_matrix/longest_increasing_path_in_matrix.go) | DFS + memoization |
| 992 | Subarrays with K Different Integers | Hard | [subarrays_with_k_different_integers/subarrays_with_k_different_integers.go](subarrays_with_k_different_integers/subarrays_with_k_different_integers.go) | Sliding window: atMost(k) - atMost(k-1) |
| 1760 | Minimum Limit of Balls in a Bag | Medium | [minimum_limit_of_balls_in_a_bag/minimum_limit_of_balls_in_a_bag.go](minimum_limit_of_balls_in_a_bag/minimum_limit_of_balls_in_a_bag.go) | Binary search on answer |
| 1235 | Maximum Profit in Job Scheduling | Hard | [maximum_profit_in_job_scheduling/maximum_profit_in_job_scheduling.go](maximum_profit_in_job_scheduling/maximum_profit_in_job_scheduling.go) | Sort by end + DP + binary search |
| 25 | Reverse Nodes in k-Group | Hard | [reverse_nodes_in_k_group/reverse_nodes_in_k_group.go](reverse_nodes_in_k_group/reverse_nodes_in_k_group.go) | Iterative reversal + recursion |
| 206 | Reverse Linked List | Easy | [reverse_linked_list/reverse_linked_list.go](reverse_linked_list/reverse_linked_list.go) | Iterative with prev pointer |
| 69 | Sqrt(x) | Easy | [sqrt_x/sqrt_x.go](sqrt_x/sqrt_x.go) | Binary search for largest m where m² ≤ x |
| 766 | Toeplitz Matrix | Easy | [toeplitz_matrix/toeplitz_matrix.go](toeplitz_matrix/toeplitz_matrix.go) | Check each cell matches top-left diagonal neighbor |
| 41 | First Missing Positive | Hard | [first_missing_positive/first_missing_positive.go](first_missing_positive/first_missing_positive.go) | Index-as-hash: place n at index n-1 |
| 33 | Search in Rotated Sorted Array | Medium | [search_in_rotated_sorted_array/search_in_rotated_sorted_array.go](search_in_rotated_sorted_array/search_in_rotated_sorted_array.go) | Binary search with sorted-half detection |
| 200 | Number of Islands | Medium | [number_of_islands/number_of_islands.go](number_of_islands/number_of_islands.go) | DFS flood fill |
| 721 | Accounts Merge | Medium | [accounts_merge/accounts_merge.go](accounts_merge/accounts_merge.go) | Union-Find on emails |
| 262 | Trips and Users | Hard | [trips_and_users/trips_and_users.sql](trips_and_users/trips_and_users.sql) | SQL — JOIN to filter banned users, GROUP BY day |
| 282 | Expression Add Operators | Hard | [expression_add_operators/expression_add_operators.go](expression_add_operators/expression_add_operators.go) | Backtracking with multiply-term tracking |

---

## Pseudo Code Summaries

### Pascal's Triangle
```
result = [[1]]
for i in 1..numRows-1:
  row = [1] + [prev[j-1]+prev[j] for j in 1..len(prev)-1] + [1]
  result.append(row)
```

### Minimum Path Sum
```
for each cell (i,j):
  dp[i][j] += min(dp[i-1][j], dp[i][j-1])   # edges handled separately
return dp[m-1][n-1]
```

### Maximum Units on a Truck
```
sort boxTypes by units desc
greedily fill truck with highest-unit boxes first
```

### Isomorphic Strings
```
maintain s->t and t->s maps
any mismatch in either direction → false
```

### Custom Sort String
```
assign rank to each char in order
sort s by rank (unranked chars get 0)
```

### Binary Tree Zigzag Level Order Traversal
```
BFS; fill level array left-to-right or right-to-left alternately
```

### Word Ladder II
```
BFS layer by layer, build parent map
DFS from endWord using parents to rebuild all shortest paths
```

### Binary Tree Right Side View
```
BFS; record last node value at each level
```

### Convert BST to Sorted Doubly Linked List
```
in-order traversal; link each node to prev; circularize head<->tail
```

### Fraction to Recurring Decimal
```
long division; track remainder→position; insert () when remainder repeats
```

### Length of Last Word
```
scan from end, skip trailing spaces, count non-space chars
```

### Single Number
```
XOR all → pairs cancel, unique survives
```

### Bus Routes
```
BFS by route; stop→routes map; count bus transfers
```

### Word Pattern II
```
backtracking; bind char→word and word→char; try all prefix splits
```

### Diagonal Traverse
```
even diagonals go up-right, odd go down-left; iterate diagonals 0..m+n-2
```

### String Compression
```
two-pointer run-length encoding in-place
```

### String to Integer (atoi)
```
skip spaces → read sign → accumulate digits → clamp to INT32 range
```

### Rank Teams by Votes
```
count votes[team][position]; sort by position votes desc, then alphabetically
```

### All O'one Data Structure
```
doubly linked list of buckets (count → set of keys); O(1) all ops
```

### Container With Most Water
```
two pointers; advance the shorter side
```

### Kth Smallest Element in a Sorted Matrix
```
binary search on value; count ≤ mid using staircase walk; shrink to k-th
```

### Palindrome Number
```
reverse second half of digits; compare with first half
```

### Palindrome Linked List
```
find mid (slow/fast) → reverse second half → compare
```

### Check If a Word Occurs As a Prefix
```
split sentence; return 1-indexed position of first word with HasPrefix
```

### Zero Array Transformation I
```
difference array to apply all decrements; prefix sum check each position
```

### Longest Increasing Path in a Matrix
```
DFS + memo; from each cell explore strictly greater neighbors
```

### Subarrays with K Different Integers
```
exactly(k) = atMost(k) - atMost(k-1); sliding window for atMost
```

### Minimum Limit of Balls in a Bag
```
binary search on max size; ops = sum(ceil(n/mid)-1)
```

### Maximum Profit in Job Scheduling
```
sort by end; dp[i] = max(dp[i-1], dp[j]+profit[i]) where j is latest non-overlapping
```

### Reverse Nodes in k-Group
```
check k nodes exist; reverse them; recurse on remainder
```

### Reverse Linked List
```
iterative prev/curr pointer swap
```

### Sqrt(x)
```
binary search for largest m where m*m <= x
```

### Toeplitz Matrix
```
each cell must equal its top-left diagonal neighbor
```

### First Missing Positive
```
cycle-sort: place n at index n-1; scan for first mismatch
```

### Search in Rotated Sorted Array
```
binary search; determine sorted half; check if target lies within it
```

### Number of Islands
```
DFS flood-fill; count connected components of '1'
```

### Accounts Merge
```
Union-Find on emails; group by root; sort each group
```

### Trips and Users (SQL)
```
JOIN to filter banned users; GROUP BY day; cancellation rate = cancelled/total
```

### Expression Add Operators
```
backtracking; track val and last-multiplied-term for * precedence
```
