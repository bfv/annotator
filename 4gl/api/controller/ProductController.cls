@controller(basePath="/api/products").
@cors(enabled="true").
class api.controller.ProductController:

  @id(type="guid").
  define property Id as char no-undo get. set.
  
  @name.
  define property Name as character no-undo get. set.


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
