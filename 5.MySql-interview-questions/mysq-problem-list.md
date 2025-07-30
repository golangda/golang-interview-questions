# ✅ 50 MySQL Interview Questions – With Answers

This document provides full answers and explanations to 50 essential MySQL interview questions. Each section includes the original question, hint, and a complete answer.

---

### 1. ❓ Write a query to find the second highest salary from an employee table.
**Hint**: Use `LIMIT` and `ORDER BY` with a subquery.
```sql
SELECT MAX(salary) AS SecondHighestSalary
FROM employees
WHERE salary < (SELECT MAX(salary) FROM employees);
```

### 2. ❓ Get the list of employees who have the same salary.
**Hint**: Use `GROUP BY salary HAVING COUNT(*) > 1`.
```sql
SELECT salary FROM employees GROUP BY salary HAVING COUNT(*) > 1;
```

### 3. ❓ Retrieve the last 5 records from a table.
**Hint**: Use `ORDER BY id DESC LIMIT 5`.
```sql
SELECT * FROM employees ORDER BY id DESC LIMIT 5;
```

### 4. ❓ Find employees who joined in the last 6 months.
**Hint**: Use `NOW()` and `INTERVAL`.
```sql
SELECT * FROM employees WHERE join_date >= NOW() - INTERVAL 6 MONTH;
```

### 5. ❓ Write a query to fetch duplicate rows from a table.
**Hint**: Use `GROUP BY` with `HAVING COUNT(*) > 1`.
```sql
SELECT column1, COUNT(*) FROM table_name GROUP BY column1 HAVING COUNT(*) > 1;
```

### 6. ❓ How to delete duplicate rows but keep one?
**Hint**: Use CTE or `ROW_NUMBER()` if available, or nested subqueries.
```sql
DELETE FROM employees
WHERE id NOT IN (
  SELECT MIN(id) FROM employees GROUP BY email
);
```

### 7. ❓ List all foreign keys in a given database.
**Hint**: Query `INFORMATION_SCHEMA.KEY_COLUMN_USAGE`.
```sql
SELECT * FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE referenced_table_name IS NOT NULL;
```

### 8. ❓ How to get the column with the maximum value in each row?
**Hint**: Use `GREATEST()` function creatively.
```sql
SELECT GREATEST(col1, col2, col3) AS max_value FROM table_name;
```

### 9. ❓ Explain and write a query to use `GROUP_CONCAT`.
**Hint**: Try `GROUP_CONCAT(column_name)`.
```sql
SELECT department_id, GROUP_CONCAT(employee_name) FROM employees GROUP BY department_id;
```

### 10. ❓ How to write a recursive query in MySQL?
**Hint**: Use `WITH RECURSIVE` for hierarchy.
```sql
WITH RECURSIVE nums AS (
  SELECT 1 AS n
  UNION ALL
  SELECT n + 1 FROM nums WHERE n < 10
)
SELECT * FROM nums;
```

### 11. ❓ Write a query to find the nth highest salary using a subquery.
**Hint**: Use correlated subquery with `LIMIT N-1,1`.
```sql
SELECT DISTINCT salary FROM employees ORDER BY salary DESC LIMIT N-1, 1;
```

### 12. ❓ Count the number of employees in each department.
**Hint**: Use `GROUP BY department_id`.
```sql
SELECT department_id, COUNT(*) FROM employees GROUP BY department_id;
```

### 13. ❓ What is the difference between WHERE and HAVING?
**Hint**: WHERE filters before aggregation, HAVING after.
- `WHERE` filters rows before grouping.
- `HAVING` filters groups after `GROUP BY`.

### 14. ❓ Use `COALESCE` to replace NULLs in a result set.
**Hint**: Use `COALESCE(column, 'default')`.
```sql
SELECT COALESCE(manager_id, 'N/A') FROM employees;
```

### 15. ❓ Write a query to transpose rows into columns.
**Hint**: Use `MAX(CASE WHEN ...)` structure.
```sql
SELECT
  employee_id,
  MAX(CASE WHEN month = 'Jan' THEN sales END) AS Jan,
  MAX(CASE WHEN month = 'Feb' THEN sales END) AS Feb
FROM sales
GROUP BY employee_id;
```
### 16. ❓ How to optimize a query using EXPLAIN?
**Hint**: Use `EXPLAIN SELECT ...`.
```sql
EXPLAIN SELECT * FROM employees WHERE department_id = 2;
```

### 17. ❓ List tables without any primary key.
**Hint**: Query `INFORMATION_SCHEMA.TABLES`.
```sql
SELECT table_name FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = 'your_db' AND table_name NOT IN (
  SELECT table_name FROM INFORMATION_SCHEMA.KEY_COLUMN_USAGE WHERE constraint_name = 'PRIMARY'
);
```

### 18. ❓ Find customers who placed more than 3 orders in the last year.
**Hint**: Use `GROUP BY customer_id HAVING COUNT(orders) > 3`.
```sql
SELECT customer_id FROM orders WHERE order_date >= CURDATE() - INTERVAL 1 YEAR GROUP BY customer_id HAVING COUNT(*) > 3;
```

### 19. ❓ Write a query using CTE (WITH clause) to sum revenue per region.
**Hint**: Use `WITH` and aggregate.
```sql
WITH regional_revenue AS (
  SELECT region, SUM(revenue) AS total FROM sales GROUP BY region
)
SELECT * FROM regional_revenue;
```

### 20. ❓ Implement pagination using LIMIT and OFFSET.
**Hint**: Use `LIMIT` with OFFSET for pages.
```sql
SELECT * FROM products LIMIT 10 OFFSET 20; -- Page 3, assuming 10 per page
```

### 21. ❓ Find the median salary from an employee table.
**Hint**: Use `ORDER BY` and `LIMIT` to simulate median.
```sql
SELECT salary FROM employees ORDER BY salary LIMIT 1 OFFSET (SELECT COUNT(*) FROM employees) DIV 2;
```

### 22. ❓ Create a view to show department-wise average salary.
**Hint**: Use `CREATE VIEW` syntax.
```sql
CREATE VIEW dept_avg_salary AS SELECT department_id, AVG(salary) FROM employees GROUP BY department_id;
```

### 23. ❓ Get the list of indexes on a table.
**Hint**: Query `SHOW INDEXES FROM tablename;`.
```sql
SHOW INDEXES FROM employees;
```

### 24. ❓ Find employees with names that start and end with vowels.
**Hint**: Use `REGEXP` with `^[aeiou].*[aeiou]$`.
```sql
SELECT * FROM employees WHERE name REGEXP '^[aeiouAEIOU].*[aeiouAEIOU]$';
```

### 25. ❓ Difference between INNER JOIN and LEFT JOIN — with example.
**Hint**: Write two sample queries to demonstrate.
```sql
-- INNER JOIN
SELECT * FROM emp e INNER JOIN dept d ON e.dept_id = d.id;
-- LEFT JOIN
SELECT * FROM emp e LEFT JOIN dept d ON e.dept_id = d.id;
```

### 26. ❓ Show rows where a timestamp is within business hours (9am–6pm).
**Hint**: Use `HOUR(timestamp_column)`.
```sql
SELECT * FROM logs WHERE HOUR(timestamp_column) BETWEEN 9 AND 17;
```

### 27. ❓ How to implement full-text search in MySQL?
**Hint**: Enable FULLTEXT index and use `MATCH ... AGAINST`.
```sql
SELECT * FROM articles WHERE MATCH(title, body) AGAINST('keyword');
```

### 28. ❓ Get all dates between two given dates using SQL.
**Hint**: Use recursive CTE or calendar table.
```sql
WITH RECURSIVE dates AS (
  SELECT '2023-01-01' AS d
  UNION ALL
  SELECT d + INTERVAL 1 DAY FROM dates WHERE d < '2023-01-10'
)
SELECT * FROM dates;
```

### 29. ❓ Delete rows from table A that don’t exist in table B (anti-join).
**Hint**: Use `NOT EXISTS` with `DELETE`.
```sql
DELETE FROM A WHERE NOT EXISTS (
  SELECT 1 FROM B WHERE A.id = B.a_id
);
```

### 30. ❓ Write a query that self joins a table (e.g. employee-manager).
**Hint**: Join table to itself using alias.
```sql
SELECT e1.name AS employee, e2.name AS manager
FROM employees e1 JOIN employees e2 ON e1.manager_id = e2.id;
```

### 31. ❓ Calculate year-over-year growth using SQL.
**Hint**: Use `LAG()` or derived tables.
```sql
SELECT year, revenue, revenue - LAG(revenue) OVER (ORDER BY year) AS growth FROM sales;
```

### 32. ❓ Detect gaps in a sequence of integers (e.g. invoice numbers).
**Hint**: Use `NOT EXISTS` and `!= previous + 1`.
```sql
SELECT t1.id + 1 AS missing_id FROM table_name t1 WHERE NOT EXISTS (
  SELECT 1 FROM table_name t2 WHERE t2.id = t1.id + 1
);
```

### 33. ❓ What’s the difference between CHAR and VARCHAR?
**Hint**: CHAR is fixed, VARCHAR is variable.
- `CHAR(n)`: Fixed-length.
- `VARCHAR(n)`: Variable-length, more efficient.

### 34. ❓ Show table size in MB including indexes.
**Hint**: Use `INFORMATION_SCHEMA.TABLES` + `DATA_LENGTH`.
```sql
SELECT table_name, ROUND((data_length + index_length)/1024/1024, 2) AS size_mb FROM INFORMATION_SCHEMA.TABLES WHERE table_schema = 'your_db';
```

### 35. ❓ How to copy a table structure without data?
**Hint**: Use `CREATE TABLE new_table LIKE old_table;`.
```sql
CREATE TABLE new_table LIKE old_table;
```

### 36. ❓ Find duplicate email addresses ignoring case sensitivity.
**Hint**: Use `LOWER()` and `GROUP BY`.
```sql
SELECT LOWER(email), COUNT(*) FROM users GROUP BY LOWER(email) HAVING COUNT(*) > 1;
```

### 37. ❓ Update records conditionally using CASE in UPDATE.
**Hint**: Use `CASE` in `SET` clause.
```sql
UPDATE employees SET bonus =
  CASE
    WHEN performance = 'A' THEN 1000
    WHEN performance = 'B' THEN 500
    ELSE 0
  END;
```

### 38. ❓ Normalize phone numbers into the same format using MySQL functions.
**Hint**: Use `REGEXP_REPLACE` or string functions.
```sql
SELECT REGEXP_REPLACE(phone, '[^0-9]', '') AS normalized FROM contacts;
```

### 39. ❓ Difference between UNION and UNION ALL — with example.
**Hint**: Write examples using same dataset.
```sql
SELECT id FROM A
UNION
SELECT id FROM B; -- Removes duplicates

SELECT id FROM A
UNION ALL
SELECT id FROM B; -- Includes duplicates
```

### 40. ❓ Get number of rows added each day for last 30 days.
**Hint**: Use `GROUP BY DATE(created_at)`.
```sql
SELECT DATE(created_at), COUNT(*) FROM logs GROUP BY DATE(created_at) ORDER BY DATE(created_at) DESC;
```

### 41. ❓ Find products not ordered in last 90 days.
**Hint**: Use `NOT IN` with subquery.
```sql
SELECT * FROM products WHERE product_id NOT IN (
  SELECT DISTINCT product_id FROM orders WHERE order_date >= NOW() - INTERVAL 90 DAY
);
```

### 42. ❓ Write a query to get the running total (cumulative sum).
**Hint**: Use `SUM() OVER(ORDER BY date)` if available.
```sql
SELECT id, date, amount, SUM(amount) OVER (ORDER BY date) AS running_total FROM transactions;
```

### 43. ❓ Get the maximum number of consecutive working days per employee.
**Hint**: Group dates and detect breaks using `DATEDIFF`.
```sql
SELECT employee_id, COUNT(*) AS streak
FROM (
  SELECT employee_id, work_date,
         DATEDIFF(work_date, ROW_NUMBER() OVER (PARTITION BY employee_id ORDER BY work_date)) AS grp
  FROM attendance
) t
GROUP BY employee_id, grp
ORDER BY streak DESC;
```

### 44. ❓ Find customers who have never placed an order.
**Hint**: Use `NOT EXISTS` with `orders` table.
```sql
SELECT * FROM customers c WHERE NOT EXISTS (
  SELECT 1 FROM orders o WHERE o.customer_id = c.id
);
```

### 45. ❓ Explain the purpose of indexing with an example.
**Hint**: Use `CREATE INDEX` and explain use.
```sql
CREATE INDEX idx_dept_id ON employees(department_id);
-- Improves performance on WHERE/ORDER BY clauses using department_id
```

### 46. ❓ Find the earliest and latest date in a table.
**Hint**: Use `MIN()` and `MAX()`.
```sql
SELECT MIN(created_at) AS earliest, MAX(created_at) AS latest FROM logs;
```

### 47. ❓ Write a query using IF and CASE to categorize age groups.
**Hint**: Use `CASE WHEN age < 18 THEN 'Minor' ...`.
```sql
SELECT name,
  CASE
    WHEN age < 18 THEN 'Minor'
    WHEN age BETWEEN 18 AND 60 THEN 'Adult'
    ELSE 'Senior'
  END AS age_group
FROM people;
```

### 48. ❓ Create a MySQL trigger that logs every delete operation.
**Hint**: Use `AFTER DELETE` trigger on a table.
```sql
CREATE TRIGGER log_deletes
AFTER DELETE ON employees
FOR EACH ROW
INSERT INTO deleted_employees_log(employee_id, deleted_at)
VALUES(OLD.id, NOW());
```

### 49. ❓ Compare performance of IN vs EXISTS.
**Hint**: Test both with large subquery datasets.
- `EXISTS` is generally faster with large subqueries.
- `IN` is simpler for static values.

### 50. ❓ Write a query to rank employees by salary within departments.
**Hint**: Use `DENSE_RANK() OVER (PARTITION BY dept ORDER BY salary DESC)`.
```sql
SELECT *, DENSE_RANK() OVER (PARTITION BY department_id ORDER BY salary DESC) AS rank
FROM employees;
```