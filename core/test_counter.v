module Counter(
  input [0:0] clk,
  input [0:0] rst,
  output [7:0] count
);

  reg [7:0] counter_reg;

  assign count = counter_reg;

endmodule
