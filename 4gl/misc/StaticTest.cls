
 /*------------------------------------------------------------------------
    File        : StaticTest
    Purpose     : 
    Syntax      : 
    Author(s)   : bronco
    Created     : Fri Jan 02 13:01:31 +01:00 2026
    Notes       : 
  ----------------------------------------------------------------------*/

block-level on error undo, throw.

class misc.StaticTest:

  constructor public StaticTest():
    
  end constructor.

  @static.
  method public static void StaticMethod():
  end method.

end class.  // StaticTest
