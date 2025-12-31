@controller(basePath="/api/products").
@cors(enabled="true").
class api.controller.ProductController:

  @get(path="/{id}").
  @response(code="200", type="Product").
  @response(code="404", type="Error").
  method public void GetProduct(input productId as integer):
  end method.

  @post.
  @requestBody(type="Product", required="true").
  method public void CreateProduct():
  end method.

  @put(path="/{id}").
  method public void UpdateProduct():
  end method.

  @delete(path="/{id}").
  @auth(role="admin").
  method public void DeleteProduct():
  end method.

end class.
