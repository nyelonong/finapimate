DROP TABLE IF EXISTS fm_user;
CREATE TABLE fm_user (
  user_id       BIGSERIAL       NOT NULL,
  email         VARCHAR(50)     NOT NULL,
  name          VARCHAR(50)     NOT NULL,
  password      VARCHAR(50)     NOT NULL,
  gender        INTEGER         NOT NULL,
  birth_date    DATE            NOT NULL,
  nik           VARCHAR(50)     NOT NULL,
  nik_valid     INTEGER         NOT NULL,
  msisdn        VARCHAR(20)     NOT NULL,
  th_amount     NUMERIC         NOT NULL,
  create_time   TIMESTAMP       NOT NULL,
  photo         VARCHAR(50),
  PRIMARY KEY (user_id)
);

DROP TABLE IF EXISTS fm_friend;
CREATE TABLE fm_friend (
  friend_id     BIGSERIAL       NOT NULL,
  user_id_a     BIGINT          NOT NULL,
  user_id_b     BIGINT          NOT NULL,
  status        INTEGER         NOT NULL,
  create_time   TIMESTAMP       NOT NULL,
  approved_time TIMESTAMP,
  PRIMARY KEY (friend_id)
);

DROP TABLE IF EXISTS fm_tx;
CREATE TABLE fm_tx (
  tx_id         BIGSERIAL       NOT NULL,
  lender_id     BIGINT          NOT NULL,
  borrower_id   BIGINT          NOT NULL,
  amount        NUMERIC         NOT NULL,
  dealine       DATE            NOT NULL,
  status        INTEGER         NOT NULL,
  create_time   TIMESTAMP       NOT NULL,
  notes         VARCHAR(512),
  PRIMARY KEY (tx_id)
);
