resource "aws_db_instance" "optical_psql" {
  identifier     = "optical-psql"
  instance_class = "db.t3.micro"
  engine         = "postgres"
  engine_version = "18.1"

  allocated_storage = 20

  db_name                     = "optical"
  username                    = "optical_admin"
  manage_master_user_password = true
  skip_final_snapshot = true

  db_subnet_group_name   = aws_db_subnet_group.optical_psql.name
  vpc_security_group_ids = [aws_security_group.optical_psql_sg.id]
}

resource "aws_db_subnet_group" "optical_psql" {
  name = "optical_psql_subnet"
  subnet_ids = [
    aws_subnet.private_a.id,
    aws_subnet.private_c.id,
  ]
}
