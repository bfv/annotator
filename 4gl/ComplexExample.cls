/* 
 * Complex example with various annotation scenarios
 */

@entity(table="customers", schema="public").
@audit(enabled="true").

// Class comment

class ComplexExample:

  /* Block comment before method */
  @rest(method="GET", path="/api/data").
  
  // Another comment
  
  @cache(ttl="300", strategy="LRU").
  method public void FetchData():
  end method.

  @deprecated.
  
  
  method public void OldMethod():
    @inmethod(type="free").
  end method.

  // Free annotation not attached to anything
  @todo(priority="high", assignee="john").

  // Another free one
  @note(text="This is just a note").

end class.
