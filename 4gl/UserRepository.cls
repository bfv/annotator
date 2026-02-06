@repository(name="UserRepository", version="1.0").
@transactional(rollback="true").
interface UserRepository:

  @query(sql="SELECT * FROM users WHERE id = ?").
  method public User GetUserById(input userId as integer).

  @query(sql="SELECT * FROM users").
  @cache(ttl="300").
  method public User[] GetAllUsers().

  @command(operation="insert").
  method public void CreateUser(input user as User).

  @command(operation="update").
  method public void UpdateUser(input user as User).

  @command(operation="delete").
  method public void DeleteUser(input userId as integer).

end interface.