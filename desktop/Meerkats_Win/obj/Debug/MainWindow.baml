   M S B A M L   `   `   `  ����G  CMeerkats_Win, Version=1.0.0.0, Culture=neutral, PublicKeyToken=null+clr-namespace:Meerkats_WinMeerkats_Win  5   	       Meerkats_Win.MainWindow   5   	   \ XPresentationFramework, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35W SPresentationCore, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35R NWindowsBase, Version=4.0.0.0, Culture=neutral, PublicKeyToken=31bf3856ad364e35D 9http://schemas.microsoft.com/winfx/2006/xaml/presentation    5   	   R NSystem.Xaml, Version=4.0.0.0, Culture=neutral, PublicKeyToken=b77a5c561934e089:x,http://schemas.microsoft.com/winfx/2006/xaml     5   	   @d2http://schemas.microsoft.com/expression/blend/2008     5   	   Jmc;http://schemas.openxmlformats.org/markup-compatibility/2006     5   	   ,localclr-namespace:Meerkats_Win     5   	     � Title$&   Meerkats File Synchronizer V1.10��5   	   $��614.4��62   $��1078.2��6A   
�������.��5	      � $��0,0,-6,-0.2q�6   $
��Left=�6    $
��1078��6;   .+�5
   
   � -   5         �Name$ header��5
      $
��Left=�6   $��98��6:   $	��Top=�6F   $
��1078��6^    � ColumnDefinitions 5      �� 5      $	��65* �6#   �� 5      $	��89* �6#   .+�5      �� -   5         ��Name$ 
conect_btn��5      $
��Sync��6)   $
��Left=�68   $��	10,42,0,0q�6S   $	��Top=�6f   $	��125��6~   $��46��6�   6�   �� 5      $H��BNotes: We will synchronize your files automatically every 20 mins ��6   $
��Left=�6a   $��
131,49,0,0q�6|   $	��Top=�6�   $	��482��6�   $��32��6�   $��15(�6�   $��1��6�   5      � 5   
   -   5        $
 File��5      $��
345,98,0,0q�6   .+�5      &   System.Windows.Controls.DataGrid  -   5          Name$ fileInfo��5      $	��487��6)   $	��Top=�66    J� RenderTransformOrigin$ -0.048,-0.096&�6N   $��0,0,0.4,-0.2q�6�   ��6t   �� *SelectedItem.Files+ �� ElementName$ directoryTreeView��6�   5      � 5   
   -   5        $ file_tree_grid��5      $
��Left=�6'   $	��487��6B   $��0,98,0,0q�6O   $	��Top=�6a   $	��345��6y   .+�5      R� -   5         R�Name$ directoryTreeView��5       � 	Resources 5      	 �� DataType"  Meerkats_Win.DirectoryRecord
 �� ItemsSource �� Orientation��%�   (��       5   /   *, +)�� 6   "	  6/   
 5      �� *Directories+.R�5      �� $ 
Horizontal=�6%   .+�5      �� ��6(   �� *	Info.Name+5      5      5      5      5       5!      5"      