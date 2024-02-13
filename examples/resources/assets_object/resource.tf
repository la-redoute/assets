resource "assets_object" "example" {
  object_type_id = "42"
  attributes_in = [
    {
      object_type_attribute_id = "42"
      object_attribute_values_in = [
        {
          value = "example"
        }
      ]
    }
  ]
}
