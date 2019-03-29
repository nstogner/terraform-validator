// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tfgcv

import (
	"context"
	"path/filepath"

	"github.com/GoogleCloudPlatform/terraform-validator/converters/google"
	"github.com/forseti-security/config-validator/pkg/api/validator"
	"github.com/forseti-security/config-validator/pkg/gcv"
	"github.com/pkg/errors"
)

// ValidateAssets instantiates GCV and audits CAI assets.
func ValidateAssets(assets []google.Asset, policyPath string) (*validator.AuditResponse, error) {
	valid, err := gcv.NewValidator(
		gcv.PolicyPath(filepath.Join(policyPath, "policies")),
		gcv.PolicyLibraryDir(filepath.Join(policyPath, "lib")),
	)
	if err != nil {
		return nil, errors.Wrap(err, "initializing gcv validator")
	}

	pbAssets := make([]*validator.Asset, len(assets))
	for i := range assets {
		pbAssets[i] = &validator.Asset{}
		if err := protoViaJSON(assets[i], pbAssets[i]); err != nil {
			return nil, errors.Wrapf(err, "converting asset %s to proto", assets[i].Name)
		}
	}

	if err := valid.AddData(&validator.AddDataRequest{
		Assets: pbAssets,
	}); err != nil {
		return nil, errors.Wrap(err, "adding data to validator")
	}

	auditResult, err := valid.Audit(context.Background())
	if err != nil {
		return nil, errors.Wrap(err, "auditing")
	}

	return auditResult, nil
}
