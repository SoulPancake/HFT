module WidthDemo(
  input [7:0] a8,
  input [15:0] b16,
  input [3:0] c4,
  output [15:0] add_out,
  output [11:0] mul_out,
  output [0:0] eq_out,
  output [27:0] cat_out
);

  assign add_out = a8 + b16;
  assign mul_out = a8 * c4;
  assign eq_out = a8 == b16;
  assign cat_out = {a8, c4, b16};

endmodule
