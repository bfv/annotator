// User service class
@service(name="UserService", version="2.0").
@deprecated(reason="Use UserServiceV2 instead").
class UserService:

  // Get user by ID
  @http(method="get", path="/users/{id}").
  @auth(required="true").
  method public void GetUser(input userId as integer):
    // Implementation here
  end method.

  @http(method="post", 
        path="/users",
        contentType="application/json").
  method public void CreateUser():
  end method.

  // This is a free annotation
  @todo(what="implement later").

  method private void HelperMethod():
  end method.

end class.
