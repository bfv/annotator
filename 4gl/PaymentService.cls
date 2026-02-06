@api(version="2.0", protocol="REST").
interface PaymentService 
  inherits BaseService:

  @http(method="POST", path="/payments").
  @validation(required="true").
  method public PaymentResult ProcessPayment(input payment as Payment).

  @http(method="GET", path="/payments/{id}").
  method public Payment GetPayment(input paymentId as integer).

  @http(method="PUT", path="/payments/{id}/refund").
  @auth(roles="admin,manager").
  method public void RefundPayment(input paymentId as integer).

  // Free annotation in interface
  @todo(priority="high", description="Add fraud detection").

end interface.