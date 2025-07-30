# 📘 50 MySQL Interview Questions with Hints

This document contains 50 carefully selected MySQL interview questions that assess both theoretical concepts and hands-on SQL proficiency. Each question is paired with a helpful hint to support your preparation.

---

| No. | Interview Problem                                                       | Hint                                                              |
| --- | ----------------------------------------------------------------------- | ----------------------------------------------------------------- |
| 1   | Write a query to find the second highest salary from an employee table. | Use `LIMIT` and `ORDER BY` with a subquery.                       |
| 2   | Get the list of employees who have the same salary.                     | Use `GROUP BY salary HAVING COUNT(*) > 1`.                        |
| 3   | Retrieve the last 5 records from a table.                               | Use `ORDER BY id DESC LIMIT 5`.                                   |
| 4   | Find employees who joined in the last 6 months.                         | Use `NOW()` and `INTERVAL`.                                       |
| 5   | Write a query to fetch duplicate rows from a table.                     | Use `GROUP BY` with `HAVING COUNT(*) > 1`.                        |
| 6   | How to delete duplicate rows but keep one?                              | Use CTE or `ROW_NUMBER()` if available, or nested subqueries.     |
| 7   | List all foreign keys in a given database.                              | Query `INFORMATION_SCHEMA.KEY_COLUMN_USAGE`.                      |
| 8   | How to get the column with the maximum value in each row?               | Use `GREATEST()` function creatively.                             |
| 9   | Explain and write a query to use `GROUP_CONCAT`.                        | Try `GROUP_CONCAT(column_name)`.                                  |
| 10  | How to write a recursive query in MySQL?                                | Use `WITH RECURSIVE` for hierarchy.                               |
| 11  | Write a query to find the nth highest salary using a subquery.          | Use correlated subquery with `LIMIT N-1,1`.                       |
| 12  | Count the number of employees in each department.                       | Use `GROUP BY department_id`.                                     |
| 13  | What is the difference between WHERE and HAVING?                        | WHERE filters before aggregation, HAVING after.                   |
| 14  | Use `COALESCE` to replace NULLs in a result set.                        | Use `COALESCE(column, 'default')`.                                |
| 15  | Write a query to transpose rows into columns.                           | Use `MAX(CASE WHEN ...)` structure.                               |
| 16  | How to optimize a query using EXPLAIN?                                  | Use `EXPLAIN SELECT ...`.                                         |
| 17  | List tables without any primary key.                                    | Query `INFORMATION_SCHEMA.TABLES`.                                |
| 18  | Find customers who placed more than 3 orders in the last year.          | Use `GROUP BY customer_id HAVING COUNT(orders) > 3`.              |
| 19  | Write a query using CTE (WITH clause) to sum revenue per region.        | Use `WITH` and aggregate.                                         |
| 20  | Implement pagination using LIMIT and OFFSET.                            | Use `LIMIT` with OFFSET for pages.                                |
| 21  | Find the median salary from an employee table.                          | Use `ORDER BY` and `LIMIT` to simulate median.                    |
| 22  | Create a view to show department-wise average salary.                   | Use `CREATE VIEW` syntax.                                         |
| 23  | Get the list of indexes on a table.                                     | Query `SHOW INDEXES FROM tablename;`.                             |
| 24  | Find employees with names that start and end with vowels.               | Use `REGEXP` with `^[aeiou].*[aeiou]$`.                           |
| 25  | Difference between INNER JOIN and LEFT JOIN — with example.             | Write two sample queries to demonstrate.                          |
| 26  | Show rows where a timestamp is within business hours (9am–6pm).         | Use `HOUR(timestamp_column)`.                                     |
| 27  | How to implement full-text search in MySQL?                             | Enable FULLTEXT index and use `MATCH ... AGAINST`.                |
| 28  | Get all dates between two given dates using SQL.                        | Use recursive CTE or calendar table.                              |
| 29  | Delete rows from table A that don’t exist in table B (anti-join).       | Use `NOT EXISTS` with `DELETE`.                                   |
| 30  | Write a query that self joins a table (e.g. employee-manager).          | Join table to itself using alias.                                 |
| 31  | Calculate year-over-year growth using SQL.                              | Use `LAG()` or derived tables.                                    |
| 32  | Detect gaps in a sequence of integers (e.g. invoice numbers).           | Use `NOT EXISTS` and `!= previous + 1`.                           |
| 33  | What’s the difference between CHAR and VARCHAR?                         | CHAR is fixed, VARCHAR is variable.                               |
| 34  | Show table size in MB including indexes.                                | Use `INFORMATION_SCHEMA.TABLES` + `DATA_LENGTH`.                  |
| 35  | How to copy a table structure without data?                             | Use `CREATE TABLE new_table LIKE old_table;`.                     |
| 36  | Find duplicate email addresses ignoring case sensitivity.               | Use `LOWER()` and `GROUP BY`.                                     |
| 37  | Update records conditionally using CASE in UPDATE.                      | Use `CASE` in `SET` clause.                                       |
| 38  | Normalize phone numbers into the same format using MySQL functions.     | Use `REGEXP_REPLACE` or string functions.                         |
| 39  | Difference between UNION and UNION ALL — with example.                  | Write examples using same dataset.                                |
| 40  | Get number of rows added each day for last 30 days.                     | Use `GROUP BY DATE(created_at)`.                                  |
| 41  | Find products not ordered in last 90 days.                              | Use `NOT IN` with subquery.                                       |
| 42  | Write a query to get the running total (cumulative sum).                | Use `SUM() OVER(ORDER BY date)` if available.                     |
| 43  | Get the maximum number of consecutive working days per employee.        | Group dates and detect breaks using `DATEDIFF`.                   |
| 44  | Find customers who have never placed an order.                          | Use `NOT EXISTS` with `orders` table.                             |
| 45  | Explain the purpose of indexing with an example.                        | Use `CREATE INDEX` and explain use.                               |
| 46  | Find the earliest and latest date in a table.                           | Use `MIN()` and `MAX()`.                                          |
| 47  | Write a query using IF and CASE to categorize age groups.               | Use `CASE WHEN age < 18 THEN 'Minor' ...`.                        |
| 48  | Create a MySQL trigger that logs every delete operation.                | Use `AFTER DELETE` trigger on a table.                            |
| 49  | Compare performance of IN vs EXISTS.                                    | Test both with large subquery datasets.                           |
| 50  | Write a query to rank employees by salary within departments.           | Use `DENSE_RANK() OVER (PARTITION BY dept ORDER BY salary DESC)`. |

---

> ✅ Use this as your master reference and practice guide for upcoming MySQL interviews. For hands-on SQL coding, explore:
>
> * [LeetCode SQL](https://leetcode.com/problemset/database/)
> * [HackerRank SQL](https://www.hackerrank.com/domains/sql)
> * [Mode Analytics SQL Tutorial](https://mode.com/sql-tutorial/)

Let me know if you’d like the solutions, explanations, or a downloadable version!
