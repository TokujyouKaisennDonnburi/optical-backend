resource "aws_security_group" "optical_api_sg" {
  name   = "optical_api_sg"
  vpc_id = aws_vpc.optical_vpc.id
  ingress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
    security_groups = [aws_security_group.optical_alb_sg.id]
  }
  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

resource "aws_security_group" "optical_psql_sg" {
  name   = "optical_psql_sg"
  vpc_id = aws_vpc.optical_vpc.id
  ingress {
    from_port       = 5432
    to_port         = 5432
    protocol        = "tcp"
    security_groups = [aws_security_group.optical_api_sg.id]
  }
}

resource "aws_security_group" "optical_redis_sg" {
  name   = "optical_redis_sg"
  vpc_id = aws_vpc.optical_vpc.id
  ingress {
    from_port       = 6379
    to_port         = 6379
    protocol        = "tcp"
    security_groups = [aws_security_group.optical_api_sg.id]
  }
}

resource "aws_security_group" "optical_alb_sg" {
  name   = "optical_alb_sg"
  vpc_id = aws_vpc.optical_vpc.id

  ingress {
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  ingress {
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

