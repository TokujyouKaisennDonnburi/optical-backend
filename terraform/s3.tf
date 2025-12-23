resource "aws_s3_bucket" "optical_images_bucket" {
  bucket              = "optical-images-bucket"
}

resource "aws_vpc_endpoint" "s3" {
  vpc_id       = aws_vpc.optical_vpc.id
  service_name = "com.amazonaws.ap-northeast-1.s3"
  vpc_endpoint_type = "Gateway"
  route_table_ids   = [
    aws_route_table.private_rt.id
  ]
}
