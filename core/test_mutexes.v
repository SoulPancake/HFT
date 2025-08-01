module MutexTest(
  input [0:0] clk,
  input [0:0] rst,
  input [0:0] priority_arb_req_0,
  input [0:0] priority_arb_req_1,
  input [0:0] priority_arb_req_2,
  input [0:0] rr_arb_req_0,
  input [0:0] rr_arb_req_1,
  output [0:0] priority_arb_grant_0,
  output [0:0] priority_arb_grant_1,
  output [0:0] priority_arb_grant_2,
  output [0:0] rr_arb_grant_0,
  output [0:0] rr_arb_grant_1
);

  reg [0:0] rr_arb_counter;

  assign priority_arb_grant_2 = priority_arb_req_2;
  assign priority_arb_grant_1 = priority_arb_req_1 && !(priority_arb_req_2);
  assign priority_arb_grant_0 = priority_arb_req_0 && !(priority_arb_req_1 || priority_arb_req_2);

  // Mutex: priority_arb (priority arbitration)
  // Request[0]: priority_arb_req_0 -> Grant[0]: priority_arb_grant_0
  // Request[1]: priority_arb_req_1 -> Grant[1]: priority_arb_grant_1
  // Request[2]: priority_arb_req_2 -> Grant[2]: priority_arb_grant_2

  // Mutex: rr_arb (round_robin arbitration)
  // Request[0]: rr_arb_req_0 -> Grant[0]: rr_arb_grant_0
  // Request[1]: rr_arb_req_1 -> Grant[1]: rr_arb_grant_1

  always @(posedge clk) begin

    rr_arb_grant_0 <= (rr_arb_counter == 0) && rr_arb_req_0;

    rr_arb_grant_1 <= (rr_arb_counter == 1) && rr_arb_req_1;

    if (rr_arb_grant_0 || rr_arb_grant_1) rr_arb_counter <= (rr_arb_counter + 1) % 2;

  end

endmodule
