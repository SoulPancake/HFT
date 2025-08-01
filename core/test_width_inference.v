module WidthInferenceTest(
  input [7:0] a,
  input [15:0] b,
  input [3:0] c,
  output [15:0] add_out,
  output [11:0] mul_out,
  output [27:0] cat_out
);

  assign add_out = a + b;
  assign mul_out = a * c;
  assign cat_out = {a, b, c};

endmodule
