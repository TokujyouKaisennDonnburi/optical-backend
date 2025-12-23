resource "aws_lb" "optical_alb" {
  name               = "optical-alb"
  internal           = false
  load_balancer_type = "application"
  security_groups    = [aws_security_group.optical_alb_sg.id]
  subnets = [
    aws_subnet.public_sub_a.id,
    aws_subnet.public_sub_c.id,
  ]
}

resource "aws_lb_target_group" "optical_alb_tg" {
  name        = "optical-alb-tg"
  port        = 8000
  protocol    = "HTTP"
  vpc_id      = aws_vpc.optical_vpc.id
  target_type = "instance"
}

resource "aws_lb_target_group_attachment" "optical_alb_tg_attach" {
  target_group_arn = aws_lb_target_group.optical_alb_tg.arn
  target_id        = aws_instance.optical_backend_api.id
  port             = 8000
}

resource "aws_acm_certificate" "optical_domain" {
  domain_name       = var.optical_domain
  validation_method = "DNS"
}

resource "aws_acm_certificate_validation" "optical_domain" {
  certificate_arn         = aws_acm_certificate.optical_domain.arn
  validation_record_fqdns = [for dvo in aws_acm_certificate.optical_domain.domain_validation_options : dvo.resource_record_name]
}

resource "aws_lb_listener" "https" {
  load_balancer_arn = aws_lb.optical_alb.arn
  port              = 443
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = aws_acm_certificate_validation.optical_domain.certificate_arn

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.optical_alb_tg.arn
  }
}
