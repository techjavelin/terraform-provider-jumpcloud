package planmodifiers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type BoolDefaultModifier struct {
	Default bool
}

func (m BoolDefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", strconv.FormatBool(m.Default))
}

func (m BoolDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", strconv.FormatBool(m.Default))
}

func (m BoolDefaultModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	if !req.AttributePlan.IsNull() {
		return
	}

	var boolval types.Bool
	diags := tfsdk.ValueAs(ctx, req.AttributePlan, boolval)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	resp.AttributePlan = types.BoolValue(m.Default)
}
