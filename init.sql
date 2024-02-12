CREATE TABLE clientes (
        id SERIAL PRIMARY KEY,
        limite INTEGER NOT NULL,
        saldo INTEGER NOT NULL
);

CREATE TABLE transacoes (
        id SERIAL PRIMARY KEY,
        cliente_id INTEGER NOT NULL,
        tipo VARCHAR(1) NOT NULL,
        valor INTEGER NOT NULL,
        descricao VARCHAR(10) NOT NULL,
        realizado_em TIMESTAMP NOT NULL DEFAULT NOW()
);

DO $$
BEGIN
INSERT INTO clientes (limite, saldo)
VALUES
    (1000 * 100, 0),
    (800 * 100, 0),
    (10000 * 100, 0),
    (100000 * 100, 0),
    (5000 * 100, 0);
END; $$