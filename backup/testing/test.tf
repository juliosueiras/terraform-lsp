variable "test_object" {
  type = object({
    test = string
  })
}

resource "aws_instance" "default" {
  test {
    test = var.test_object.
  }
}
