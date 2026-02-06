/* Interface example with contract annotations */

@contract(type="service").
interface INotificationService:

  @async(timeout="5000").
  method public void SendEmail(input recipient as character, input message as character).
  
  @async(timeout="3000").
  method public void SendSMS(input phone as character, input message as character).

  // Free annotation in interface
  @todo(priority="high", description="Add retry mechanism").

end interface.