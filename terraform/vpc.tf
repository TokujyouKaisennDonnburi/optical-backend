resource "aws_vpc" "optical_vpc" {
  tags = {
    Name = "optical-vpc"
  }
  cidr_block           = "10.0.0.0/16"
  enable_dns_support   = true
  enable_dns_hostnames = true
}

resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.optical_vpc.id
}

resource "aws_subnet" "public_sub_a" {
  vpc_id                  = aws_vpc.optical_vpc.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "ap-northeast-1a"
  map_public_ip_on_launch = true
}

resource "aws_subnet" "public_sub_c" {
  vpc_id                  = aws_vpc.optical_vpc.id
  cidr_block              = "10.0.2.0/24"
  availability_zone       = "ap-northeast-1c"
  map_public_ip_on_launch = true
}

resource "aws_route_table" "public_rt" {
  vpc_id = aws_vpc.optical_vpc.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }
}

resource "aws_route_table" "private_rt" {
  vpc_id = aws_vpc.optical_vpc.id
}

resource "aws_route_table_association" "public_assoc_a" {
  subnet_id      = aws_subnet.public_sub_a.id
  route_table_id = aws_route_table.public_rt.id
}

resource "aws_route_table_association" "public_assoc_c" {
  subnet_id      = aws_subnet.public_sub_c.id
  route_table_id = aws_route_table.public_rt.id
}

resource "aws_subnet" "private_a" {
  vpc_id            = aws_vpc.optical_vpc.id
  cidr_block        = "10.0.10.0/24"
  availability_zone = "ap-northeast-1a"
}

resource "aws_subnet" "private_c" {
  vpc_id            = aws_vpc.optical_vpc.id
  cidr_block        = "10.0.11.0/24"
  availability_zone = "ap-northeast-1c"
}
