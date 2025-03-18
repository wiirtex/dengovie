package postgres

import (
	"context"
	"dengovie/internal/store/types"
	"fmt"
	"strings"
)

func (r *Repo) ListUserDebts(ctx context.Context, input types.ListUserDebtsInput) ([]types.UserDebt, error) {

	query := `
select d.another_user_id, u.name, sum(d.direction * d.amount)
from users u
inner join debts d on u.id = d.user_id
where d.user_id = $1
group by d.another_user_id, u.id
`

	rows, err := r.db.QueryContext(ctx, query, input.UserID)
	if err != nil {
		return nil, fmt.Errorf("db.QueryContext: %w", err)
	}

	debts := make([]types.UserDebt, 0, 5)
	for rows.Next() {
		var debt types.UserDebt
		errScan := rows.Scan(&debt.AnotherUser.ID, &debt.AnotherUser.Name, &debt.Amount)
		if errScan != nil {
			return nil, fmt.Errorf("rows.Scan: %w", errScan)
		}
		fmt.Println("add debt", debt)

		debts = append(debts, debt)
	}

	return debts, nil
}

func (r *Repo) CreateEmptyDebts(ctx context.Context, input types.CreateEmptyDebtsInput) error {
	values := strings.Builder{}
	args := make([]interface{}, 0, 2*len(input.AnotherUserIDs))

	for i := 0; i/2 < len(input.AnotherUserIDs); i += 2 {
		values.WriteString(fmt.Sprintf("($%d, $%d, 1, 0),", i+1, i+2))
		values.WriteString(fmt.Sprintf("($%d, $%d, -1, 0),", i+2, i+1))
		args = append(args, input.UserID, input.AnotherUserIDs[i/2])
	}
	valuesStr := values.String()
	fmt.Println(valuesStr, args, input.AnotherUserIDs)

	query := fmt.Sprintf(`
insert into debts (user_id, another_user_id, direction, amount) 
values %s
on conflict (user_id, another_user_id, direction) do nothing`, valuesStr[:len(valuesStr)-1])

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	fmt.Println(res.RowsAffected())

	return nil
}

func (r *Repo) ShareDebt(ctx context.Context, input types.ShareDebtInput) error {

	values := strings.Builder{}
	args := make([]interface{}, 0, len(input.ChangeDebtAmount))

	for i := 0; i/3 < len(input.ChangeDebtAmount); i += 3 {
		values.WriteString(fmt.Sprintf("(1, $%d::bigint, $%d::bigint, $%d::bigint), (-1, $%d, $%d, $%d),", i+1, i+2, i+3, i+1, i+3, i+2))
		args = append(args, input.ChangeDebtAmount[i/3].Amount, int64(input.UserID), int64(input.ChangeDebtAmount[i/3].UserID))
	}

	valuesStr := values.String()
	fmt.Println(valuesStr)
	fmt.Printf("%#v\n", args)
	for _, arg := range args {
		fmt.Printf("%#v, %T\n", arg, arg)
	}

	query := fmt.Sprintf(`
update debts d
set amount = d.amount + nv.amount_diff
from
    (
        values
        %s
    ) as nv (direction, amount_diff, user_id, another_user_id)
where d.direction = nv.direction and d.user_id = nv.user_id and d.another_user_id = nv.another_user_id  
	`, valuesStr[:len(valuesStr)-1])

	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("db.ExecContext: %w", err)
	}

	fmt.Println(res.RowsAffected())
	return nil
}
