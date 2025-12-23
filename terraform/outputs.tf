output "vpc_id" {
  value       = aws_vpc.optical_vpc.id
  description = "vpc-id"
}

output "alb_dns_name" {
  description = "ALB dns name"
  value       = aws_lb.optical_alb.dns_name
}

output "acm_domain_validation" {
  description = "domain validation"
  value = aws_acm_certificate.optical_domain.domain_validation_options
}

output "s3_bucket_domain_name" {
  description = "The bucket domain name"
  value       = aws_s3_bucket.optical_images_bucket.bucket_domain_name
}

output "s3_bucket_arn" {
  value       = aws_s3_bucket.optical_images_bucket.arn
  description = "s3 bucket arn"
}
