// Copyright 2024 The Score Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package command

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/score-spec/score-implementation-sample/internal/version"
)

const (
	versionCmdFileNoLogo = "no-logo"
	versionCmdFileNoUpdateCheck = "no-update-check"
	logo =`
                   ...    ..-----------           
               .......   .....--------            
           .........     ......------             
       .........        .....                     
      .......           ....   ..                 
       ..........      .....   .....-             
           ........   ......   .......---         
              .....   .....       ......----      
                      ....          .....---      
              ...........       .........         
           ..............    ........-            
         ................   ......                
         ...............    ..                    
          -...........
	`
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version for score-implementation-sample and new version to update if available.",
	Args:  cobra.NoArgs,
	CompletionOptions: cobra.CompletionOptions{
		HiddenDefaultCmd: true,
	},
	SilenceErrors: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true

		if v, _ := cmd.Flags().GetBool(versionCmdFileNoLogo); !v {
			fmt.Println(logo)
		}
		fmt.Println("score-implementation-sample", version.BuildVersionString())

		return nil
	},
}

func init() {
	versionCmd.Flags().Bool(versionCmdFileNoLogo, false, "Do not show the Score logo")
	rootCmd.AddCommand(versionCmd)
}
