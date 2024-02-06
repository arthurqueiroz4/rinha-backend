CREATE TABLE clientes (
  id SERIAL PRIMARY KEY,
  limite INTEGER NOT NULL,
  saldo INTEGER NOT NULL
);

CREATE TABLE transacao (
  id SERIAL PRIMARY KEY,
  cliente_id INTEGER NOT NULL,
  tipo VARCHAR(1) NOT NULL,
  valor INTEGER NOT NULL,
  descricao VARCHAR(10) NOT NULL,
  realizado TIMESTAMP NOT NULL DEFALT NOW()
);

CREATE INDEX idx_extrato ON transacoes (id DESC);

DO $$
BEGIN
  INSERT INTO clientes (nome, limite)
  VALUES
    ('o barato sai caro', 1000 * 100),
    ('zan corp ltda', 800 * 100),
    ('les cruders', 10000 * 100),
    ('padaria joia de cocaia', 100000 * 100),
    ('kid mais', 5000 * 100);
END; $$