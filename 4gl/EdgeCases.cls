@test(value="with ~"quotes~" escaped", other="normal").
class EdgeCases:

  @attr(empty="").
  method public void EmptyValue():
  end.

  @simple.
  method public void NoAttributes():
  end.

  @multiline(
    first="value1",
    second="value2",
    third="value3"
  ).
  method public void MultilineAnnotation():
  end.

  /* This annotation is in a comment and should be ignored
  @ignored(should="not appear").
  */
  method public void CommentedOut():
  end.

  // @alsoignored(in="comment").
  method public void AlsoIgnored():
  end.

  method public void NoAnnotation():
    
    var int i, j.
    
    // The following line has an invalid annotation inside a comment
    assign /* blabla
              @RETURN is not a valid identifier */
      i = 1
      j = 2
      .

  end.

end class.
