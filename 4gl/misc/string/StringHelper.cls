@utility.
class misc.string.StringHelper:

  @pure.
  method public static character Trim(input str as character):
    return trim(str).
  end method.

  @pure.
  @cached.
  method public static character Upper(input str as character):
    return caps(str).
  end method.

end class.
