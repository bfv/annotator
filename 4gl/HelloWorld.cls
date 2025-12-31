@todo(version="13.0", what="refactor generics").
class HelloWorld:

  @http(method="get", comment="Returns greeting").
  method public void GetHello():
    message "Hello World".
  end method.

  @http(method="post").
  @validate(required="true").
  method public void PostHello():
    message "Posted".
  end method.

end class.
