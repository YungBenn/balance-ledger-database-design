-- Money in (deposit, refund) is positive, money out (withdrawal, purchase) is negative.
ALTER TABLE ledger
    ADD CONSTRAINT ledger_add_sign_check CHECK (
        (type IN ('deposit', 'refund') AND add > 0)
        OR (type IN ('withdrawal', 'purchase') AND add < 0)
    );
